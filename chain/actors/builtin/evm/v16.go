package evm

import (
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/abi"
	evm16 "github.com/filecoin-project/go-state-types/builtin/v16/evm"

	"github.com/filecoin-project/lotus/chain/actors/adt"
)

var _ State = (*state16)(nil)

func load16(store adt.Store, root cid.Cid) (State, error) {
	out := state16{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make16(store adt.Store, bytecode cid.Cid) (State, error) {
	out := state16{store: store}
	s, err := evm16.ConstructState(store, bytecode)
	if err != nil {
		return nil, err
	}

	out.State = *s

	return &out, nil
}

type state16 struct {
	evm16.State
	store adt.Store
}

func (s *state16) Nonce() (uint64, error) {
	return s.State.Nonce, nil
}

func (s *state16) IsAlive() (bool, error) {
	return s.State.Tombstone == nil, nil
}

func (s *state16) GetState() interface{} {
	return &s.State
}

func (s *state16) GetBytecodeCID() (cid.Cid, error) {
	return s.State.Bytecode, nil
}

func (s *state16) GetBytecodeHash() ([32]byte, error) {
	return s.State.BytecodeHash, nil
}

func (s *state16) GetBytecode() ([]byte, error) {
	bc, err := s.GetBytecodeCID()
	if err != nil {
		return nil, err
	}

	var byteCode abi.CborBytesTransparent
	if err := s.store.Get(s.store.Context(), bc, &byteCode); err != nil {
		return nil, err
	}

	return byteCode, nil
}
