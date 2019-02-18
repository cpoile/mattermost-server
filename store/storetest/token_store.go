// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenStore(t *testing.T, ss store.Store) {
	t.Run("", func(t *testing.T) { testTokenStore(t, ss) })
	t.Run("PermanentDeleteBatch", func(t *testing.T) { testAuditStorePermanentDeleteBatch(t, ss) })
}

func testTokenStore(t *testing.T, ss store.Store) {
	passwordRecoveryToken := model.NewToken("password_recovery", "foobar")
	store.Must(ss.Token().Save(passwordRecoveryToken))

	r := <-ss.Token().GetByToken(passwordRecoveryToken.Token)
	token := r.Data.(*model.Token)
	assert.NotNil(t, token)
	assert.Equal(t, passwordRecoveryToken.Extra, token.Extra)

	r2 := <-ss.Token().GetByTypeAndExtra("password_recovery", passwordRecoveryToken.Extra)
	tokens := r2.Data.([]*model.Token)
	assert.Len(t, tokens, 1)
	assert.Equal(t, passwordRecoveryToken.Extra, tokens[0].Extra)

	passwordRecoveryToken2 := model.NewToken("password_recovery", "foobar")
	store.Must(ss.Token().Save(passwordRecoveryToken2))

	r3 := <-ss.Token().GetByTypeAndExtra("password_recovery", passwordRecoveryToken.Extra)
	tokens = r3.Data.([]*model.Token)
	assert.Len(t, tokens, 2)
	assert.Nil(t, r3.Err)

	store.Must(ss.Token().Delete(passwordRecoveryToken.Token))
	r4 := <-ss.Token().GetByToken(passwordRecoveryToken.Token)
	assert.NotNil(t, r4.Err)
}
