// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package accounts implements high level Ethereum account management.
package accounts

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Account represents an Ethereum account located at a specific location defined
// by the optional URL field.
// Account 구조체는 이더리움 계정을 나타내는데, 이는 옵션으로 URL field로 정의된 특정 위치에 위치할 수 있다.
type Account struct {
	Address common.Address `json:"address"` // Ethereum account address derived from the key
	URL     URL            `json:"url"`     // Optional resource locator within a backend
}

// Wallet represents a software or hardware wallet that might contain one or more
// accounts (derived from the same seed).
// Wallet 인터페이스는 동일한 seed로부터 유도된 하나 혹은 그 이상의 계정을 포함하는 소프트웨어
// 혹은 하드웨어 wallet을 나타낸다.
type Wallet interface {
	// URL retrieves the canonical path under which this wallet is reachable. It is
	// user by upper layers to define a sorting order over all wallets from multiple
	// backends.
	URL() URL

	// Status returns a textual status to aid the user in the current state of the
	// wallet. It also returns an error indicating any failure the wallet might have
	// encountered.
	// Status 메서드는 현재 wallet의 상태를 유저에게 textual status를 반환한다.
	// 또한, wallet이 이전에 발생한 error들을 반환하여 준다.
	Status() (string, error)

	// Open initializes access to a wallet instance. It is not meant to unlock or
	// decrypt account keys, rather simply to establish a connection to hardware
	// wallets and/or to access derivation seeds.

	// Open 메서드는 wallet 지갑 객체에 접근하기 위한  초기화를 진행한다. 실제로 언락 혹은
	// 어카운트 키를 복호화하는 것이 아니라 단순히 하드웨어 지갑에 대한 연결을 설정하거나
	// 시드에에 접근하는 것과 같은 단순한 작업이다.

	// The passphrase parameter may or may not be used by the implementation of a
	// particular wallet instance. The reason there is no passwordless open method
	// is to strive towards a uniform wallet handling, oblivious to the different
	// backend providers.

	// Passphrase 파라미터는 특정한 지갑의 경우에는 사용되지 않을 수도 있다.
	// Passphrase를 넣지 않는 open 메서드가 없는 이유는 다른 백엔드 제공자와 무관하게
	// 지갑 관련하여 공통화된 방법을 제공하기 위함입니다.

	// Please note, if you open a wallet, you must close it to release any allocated
	// resources (especially important when working with hardware wallets).

	// 만약 지갑을 open 했으면 (특히 하드웨어 지갑과 연동했을 때에는 더욱 중요한데) 리소스 해제를
	// 위하여 반드시 close 해야한다.

	Open(passphrase string) error

	// Close releases any resources held by an open wallet instance.
	// Open 지갑 객체에 의하여 잡혀진 모든 리소스를 해제한다.
	Close() error

	// Accounts retrieves the list of signing accounts the wallet is currently aware
	// of. For hierarchical deterministic wallets, the list will not be exhaustive,
	// rather only contain the accounts explicitly pinned during account derivation.

	// 계정 함수는 현재 지갑이 알고 있는 signing accounts의 리스트를 반환한다.
	// HD (계층적으로 주소를 생성하는) 지갑의 경우에는
	Accounts() []Account

	// Contains returns whether an account is part of this particular wallet or not.
	// Contains 메소드는 파라미터인 account가 해당 특정 지갑의 일부인지 아닌지를 반환해준다.
	Contains(account Account) bool

	// Derive attempts to explicitly derive a hierarchical deterministic account at
	// the specified derivation path. If requested, the derived account will be added
	// to the wallet's tracked account list.

	// Devier 메소드는 특정한 derivation path에 대해서 계층 유도 계정을 유도한다.
	// 요청이 주어지면 유도 계정은 지갑의 계정 리스트에 추가된다.
	Derive(path DerivationPath, pin bool) (Account, error)

	// SelfDerive sets a base account derivation path from which the wallet attempts
	// to discover non zero accounts and automatically add them to list of tracked
	// accounts.

	// SelfDrive 메소드는 지갑이 non zero 계정을 찾을려고 하는 지갑 기본 유도 패스의 기본을
	// 설정하며 이는 자동으로 계정 리스트에 추가된다.

	// Note, self derivaton will increment the last component of the specified path
	// opposed to decending into a child path to allow discovering accounts starting
	// from non zero components.

	// You can disable automatic account discovery by calling SelfDerive with a nil
	// chain state reader.
	SelfDerive(base DerivationPath, chain ethereum.ChainStateReader)

	// SignHash requests the wallet to sign the given hash.
	// SignHash 메소드는 지갑에서 주어진 Hash에 대해서 서명하도록 요청한다.

	// It looks up the account specified either solely via its address contained within,
	// or optionally with the aid of any location metadata from the embedded URL field.

	//

	// If the wallet requires additional authentication to sign the request (e.g.
	// a password to decrypt the account, or a PIN code o verify the transaction),
	// an AuthNeededError instance will be returned, containing infos for the user
	// about which fields or actions are needed. The user may retry by providing
	// the needed details via SignHashWithPassphrase, or by other means (e.g. unlock
	// the account in a keystore).
	SignHash(account Account, hash []byte) ([]byte, error)

	// SignTx requests the wallet to sign the given transaction.
	// SignTx 메소드는 지갑에 주어진 트랜잭션에 서명하도록

	// It looks up the account specified either solely via its address contained within,
	// or optionally with the aid of any location metadata from the embedded URL field.
	//
	// If the wallet requires additional authentication to sign the request (e.g.
	// a password to decrypt the account, or a PIN code o verify the transaction),
	// an AuthNeededError instance will be returned, containing infos for the user
	// about which fields or actions are needed. The user may retry by providing
	// the needed details via SignTxWithPassphrase, or by other means (e.g. unlock
	// the account in a keystore).
	SignTx(account Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	// SignHashWithPassphrase requests the wallet to sign the given hash with the
	// given passphrase as extra authentication information.
	//
	// It looks up the account specified either solely via its address contained within,
	// or optionally with the aid of any location metadata from the embedded URL field.
	SignHashWithPassphrase(account Account, passphrase string, hash []byte) ([]byte, error)

	// SignTxWithPassphrase requests the wallet to sign the given transaction, with the
	// given passphrase as extra authentication information.
	//
	// It looks up the account specified either solely via its address contained within,
	// or optionally with the aid of any location metadata from the embedded URL field.
	SignTxWithPassphrase(account Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}

// Backend is a "wallet provider" that may contain a batch of accounts they can
// sign transactions with and upon request, do so.
type Backend interface {
	// Wallets retrieves the list of wallets the backend is currently aware of.
	//
	// The returned wallets are not opened by default. For software HD wallets this
	// means that no base seeds are decrypted, and for hardware wallets that no actual
	// connection is established.
	//
	// The resulting wallet list will be sorted alphabetically based on its internal
	// URL assigned by the backend. Since wallets (especially hardware) may come and
	// go, the same wallet might appear at a different positions in the list during
	// subsequent retrievals.
	Wallets() []Wallet

	// Subscribe creates an async subscription to receive notifications when the
	// backend detects the arrival or departure of a wallet.
	Subscribe(sink chan<- WalletEvent) event.Subscription
}

// WalletEventType represents the different event types that can be fired by
// the wallet subscription subsystem.
type WalletEventType int

const (
	// WalletArrived is fired when a new wallet is detected either via USB or via
	// a filesystem event in the keystore.
	WalletArrived WalletEventType = iota

	// WalletOpened is fired when a wallet is successfully opened with the purpose
	// of starting any background processes such as automatic key derivation.
	WalletOpened

	// WalletDropped
	WalletDropped
)

// WalletEvent is an event fired by an account backend when a wallet arrival or
// departure is detected.
type WalletEvent struct {
	Wallet Wallet          // Wallet instance arrived or departed
	Kind   WalletEventType // Event type that happened in the system
}
