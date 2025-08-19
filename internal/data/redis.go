package data

import (
	"context"
	"fmt"
)

func userInfoKey(uuid string) string {
	return fmt.Sprintf("userInfo:%s", uuid)
}

func (r *userRepo) GetUserInfo(ctx context.Context, uuid string) (info string, err error) {
	return r.data.rdb.Get(ctx, userInfoKey(uuid)).String(), nil
}
