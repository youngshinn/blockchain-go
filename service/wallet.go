package service

import (
	"block-test/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hacpy/go-ethereum/common/hexutil"
)

//func (s *Service) newKeyPair() (string, string, error) {
//
//	p256 := elliptic.P256() //암호화 방식 (타원 곡선 암호화)
//	// ECSDA 개인 키 생성(private)
//	if private, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
//		return "", "", err // 키 생성 실패시 오류 반환
//	} else if private == nil {
//		return "", "", errors.New("pk is nil") // private == nil
//	} else {
//		privatekeyBytes := crypto.FromECDSA(private)  //개인 private 키 값을 byte값으로 변환
//		privateKey := hexutil.Encode(privatekeyBytes) // 문자열로 인코딩 (이더리움 서명 표준에 맞춤)
//
//		againPrivateKey, err := crypto.HexToECDSA(privateKey)
//		if err != nil {
//			msg := fmt.Sprintf("ECDSA private key is invalid %v", err.Error())
//			panic(msg)
//		} // 키 값에서 0x 값을 제거 후 반환
//
//		publicKey := againPrivateKey.Public() // 변환된 개인 키를 통해 공개 키 가져오기
//		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
//		if !ok {
//			msg := fmt.Sprintf("public key is not of type *ecdsa.PublicKey")
//			panic(msg)
//		}
//		//공개 키를 이더리움 주소로 변환
//		address := crypto.PubkeyToAddress(*publicKeyECDSA)
//
//		return privateKey, hexutil.Encode(address[:]), nil
//
//	}
//}

func (s *Service) newKeyPair() (string, string, error) {
	// 타원 곡선 P-256 사용
	p256 := elliptic.P256()

	// ECDSA 개인 키 생성
	private, err := ecdsa.GenerateKey(p256, rand.Reader)
	if err != nil {
		return "", "", err
	}
	if private == nil {
		return "", "", errors.New("private key generation failed (nil)")
	}

	// 개인 키를 바이트 배열로 변환 후 Hex 인코딩
	privateKeyBytes := crypto.FromECDSA(private)
	privateKeyHex := hexutil.Encode(privateKeyBytes) // "0x" 포함된 16진수 문자열

	// "0x" 제거
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	// 공개 키 변환
	publicKey, ok := private.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("failed to convert public key to ECDSA format")
	}

	// 공개 키를 이더리움 주소로 변환
	address := crypto.PubkeyToAddress(*publicKey)

	// 개인 키와 이더리움 주소 반환
	return privateKeyHex, address.Hex(), nil
}

func (s *Service) MakeWallet() *types.Wallet {
	wallet := types.Wallet{
		Balance: "0",
	}
	var err error

	if wallet.PrivateKey, wallet.PublicKey, err = s.newKeyPair(); err != nil {
		return nil
	} else if err = s.repository.CreateNewWallet(&wallet); err != nil {
		return nil
	} else {
		return &wallet
	}
}

func (s *Service) GetWallet(pk string) (*types.Wallet, error) {

	if wallet, err := s.repository.GetWallet(pk); err != nil {
		return nil, err
	} else {
		return wallet, nil
	}
}

//func (s *Service) GetWalletByKey(publickey string) (*types.Wallet, error) {
//	if wallet, err := s.repository.GetWalletByPublickey(publickey); err != nil {
//		return nil, err
//	} else {
//		return wallet, nil
//	}
//}
