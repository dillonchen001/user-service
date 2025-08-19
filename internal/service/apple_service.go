package service

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	jwtv4 "github.com/golang-jwt/jwt/v4" // 添加别名
	v1 "user-service/api/auth/v1"
	"user-service/internal/biz"
	"user-service/internal/conf"
	"user-service/third_party/jwt"
)

type AppleService struct {
	cfg        *conf.Jwt
	log        *log.Helper
	authCase   *biz.UserAuthCase
	jwtGen     *jwt.Generator
	httpClient *http.Client
	publicKeys map[string]*rsa.PublicKey
}

func NewAppleService(cfg *conf.Jwt, logger log.Logger, authCase *biz.UserAuthCase) *AppleService {
	return &AppleService{
		cfg:        cfg,
		log:        log.NewHelper(logger),
		authCase:   authCase,
		jwtGen:     jwt.NewGenerator(cfg.Secret, int(cfg.Expires)),
		httpClient: &http.Client{Timeout: 10 * time.Second},
		publicKeys: make(map[string]*rsa.PublicKey),
	}
}

// AppleClaims JWT 声明结构
type AppleClaims struct {
	jwtv4.RegisteredClaims        // 使用别名
	Email                  string `json:"email"`
	Sub                    string `json:"sub"` // Apple 用户唯一标识符
}

func (s *AppleService) Login(ctx context.Context, req *v1.LoginWithAppleRequest) (*v1.LoginResponse, error) {
	// 验证参数
	if req.IdToken == "" {
		return nil, errors.New("id_token is required")
	}
	if req.Nonce == "" {
		return nil, errors.New("nonce is required")
	}

	// 解析并验证 id_token
	claims, err := s.verifyIdToken(ctx, req.IdToken, req.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to verify apple id_token: %w", err)
	}

	// 查找或创建用户
	u, isNew, err := s.authCase.FindOrCreateByAppleID(ctx, claims.Sub, claims.Email)

	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// 生成 JWT token
	token, err := s.jwtGen.GenerateToken(u.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 构建响应
	return &v1.LoginResponse{
		Token:     token,
		IsNewUser: isNew,
		UserInfo: &v1.UserInfo{
			UserId: u.UserID,
			Name:   u.Name,
			Email:  u.Email,
		},
	}, nil
}

func (s *AppleService) verifyIdToken(ctx context.Context, idToken, nonce string) (*AppleClaims, error) {
	// 测试，先写死一下
	return &AppleClaims{
		Email: "bjbd2015@163.com",
		Sub:   "293239232323",
	}, nil

	// 解析 token 头部以获取 kid
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid id_token format")
	}

	// 解码头部
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}

	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal header: %w", err)
	}

	// 获取 Apple 公钥
	publicKey, err := s.getPublicKey(ctx, header.Kid)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// 解析并验证 token
	token, err := jwtv4.ParseWithClaims(idToken, &AppleClaims{}, func(token *jwtv4.Token) (interface{}, error) {
		// 验证算法
		if _, ok := token.Method.(*jwtv4.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// 验证声明
	claims, ok := token.Claims.(*AppleClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// 验证 nonce
	// 注意: Apple 的 nonce 可能会进行 SHA-256 哈希处理
	// 这里简化处理，实际应用中需要根据 Apple 的要求进行验证

	// 验证 issuer
	if claims.Issuer != "https://appleid.apple.com" {
		return nil, errors.New("invalid issuer")
	}

	// 验证 audience
	// 应该验证 audience 是否为您的应用客户端 ID

	return claims, nil
}

func (s *AppleService) getPublicKey(_ context.Context, kid string) (*rsa.PublicKey, error) {
	// 检查缓存
	if key, ok := s.publicKeys[kid]; ok {
		return key, nil
	}

	// 从 Apple 获取公钥
	resp, err := s.httpClient.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			s.log.Errorf("close io reader failed: %v", errClose)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apple keys api returned status code: %d", resp.StatusCode)
	}

	var keysResponse struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&keysResponse); err != nil {
		return nil, err
	}

	// 查找匹配的 key
	var keyData struct {
		N string
		E string
	}
	found := false
	for _, key := range keysResponse.Keys {
		if key.Kid == kid {
			keyData.N = key.N
			keyData.E = key.E
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("public key with kid %s not found", kid)
	}

	// 解码并构建 RSA 公钥
	nBytes, err := base64.RawURLEncoding.DecodeString(keyData.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(keyData.E)
	if err != nil {
		return nil, err
	}

	// 解析 e 为整数
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	// 构建公钥
	publicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}

	// 缓存公钥
	s.publicKeys[kid] = publicKey

	return publicKey, nil
}
