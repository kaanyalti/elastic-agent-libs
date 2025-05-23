// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build go1.24 && requirefips

package keystore

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldRaiseAndErrorWhenVersionDontMatch(t *testing.T) {
	temporaryPath := GetTemporaryKeystoreFile(t)
	defer os.Remove(temporaryPath)

	badVersion := `v1pqH8nRJNCuKLrAHwATQuHpdLcP84sATrxtKMWTvapZTRcoEODVJKf2dsHXiOhSMh1EFrJTikON2oF5wZv4IM37lkJ6wt79MCFaXDqlNxBQtIA9w6vaxWnbS+92rQqtka7WrzTxal1Pd3mcK0o+ow7EAJg553UvxBqA==`

	f, err := os.OpenFile(temporaryPath, os.O_CREATE|os.O_WRONLY, 0600)
	require.NoError(t, err)
	_, _ = f.WriteString(badVersion)
	err = f.Close()
	require.NoError(t, err)

	_, err = NewFileKeystoreWithPassword(temporaryPath, NewSecureString([]byte("")))
	if assert.Error(t, err, "Expect version check error") {
		assert.Equal(t, err, fmt.Errorf("keystore format doesn't match expected version: 'v2' got 'v1'"))
	}
}

func TestOpensV2(t *testing.T) {
	ks, err := NewFileKeystoreWithPassword(filepath.Join("testdata", "keystore.v2"), NewSecureString([]byte("")))
	require.NoError(t, err)
	ls, err := AsListingKeystore(ks)
	require.NoError(t, err)
	keys, err := ls.List()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, keys[0], "key")
}

func TestFailsToOpenV1(t *testing.T) {
	_, err := NewFileKeystoreWithPassword(filepath.Join("testdata", "keystore.v1"), NewSecureString([]byte("")))
	require.Error(t, err)
}
