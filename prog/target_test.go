// Copyright 2024 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package prog

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredGlobs(t *testing.T) {
	assert.Equal(t, requiredGlobs("aa/bb"), []string{"aa/bb"})
	assert.Equal(t, requiredGlobs("aa:bb"), []string{"aa", "bb"})
	assert.Equal(t, requiredGlobs("aa:bb:-cc:dd"), []string{"aa", "bb", "dd"})
}

func TestPopulateGlob(t *testing.T) {
	assert.Empty(t, populateGlob("aa", map[string][]string{
		"bb": {"c"},
	}))
	assert.Equal(t, []string{"d", "e"}, populateGlob("aa", map[string][]string{
		"aa": {"d", "e"},
		"bb": {"c"},
	}))
	assert.Equal(t, []string{"d", "e", "f"}, populateGlob("aa:cc", map[string][]string{
		"aa": {"d", "e"},
		"bb": {"c"},
		"cc": {"f", "d"},
	}))
	assert.Equal(t, []string{"d", "f"}, populateGlob("aa:cc:-e", map[string][]string{
		"aa": {"d", "e"},
		"bb": {"c"},
		"cc": {"f", "d"},
	}))
}

func TestEmptyGlobGenerateAndMutate(t *testing.T) {
	target, err := GetTarget("linux", "amd64")
	assert.NoError(t, err)
	r := newRand(target, rand.NewSource(0))
	s := newState(target, nil, nil)
	typ := &BufferType{
		TypeCommon: TypeCommon{
			IsVarlen: true,
		},
		Kind: BufferGlob,
	}
	typ.setRef(1)

	generated, _ := typ.generate(r, s, DirIn)
	assert.Empty(t, generated.(*DataArg).Data())

	arg := MakeDataArg(typ, DirIn, []byte("/tmp/not-a-glob-match"))
	_, retry, _ := typ.mutate(r, s, arg, ArgCtx{})
	assert.False(t, retry)
	assert.Empty(t, arg.Data())
}
