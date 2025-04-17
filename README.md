

#  Go로 구현한 블록체인 시스템 (Custom Blockchain in Go)

Go 언어로 직접 구현한 블록체인 애플리케이션입니다.  
SHA-256 기반 블록 연결, POW(Proof of Work), 머클 트리, 지갑 주소 생성 등  
블록체인의 핵심 메커니즘을 직접 구현하여 구조와 원리를 체득하는 것을 목표로 합니다.

---

##  프로젝트 개요

이 프로젝트는 블록체인의 구조를 Go 언어 기반으로 직접 설계하고 구현한 학습 및 실전용 프로젝트입니다.  
블록 생성, 작업 증명(Proof of Work), 머클 트리, 지갑 생성 등 핵심 요소들을 포함하며,  
중앙 서버 없이 데이터 무결성과 분산 합의를 체험할 수 있도록 구성되어 있습니다.

---

##  사용 기술 스택

| 분류        | 내용 |
|-------------|------|
| Language    | Go (Golang) |
| 암호화 알고리즘 | SHA-256, ECDSA |
| 자료구조     | Merkle Tree |
| 설정 관리   | TOML (`environment.toml`) |
| 구조 설계   | Go 패키지 모듈화 (`service`, `repository`, `types`) |

---

##  폴더 구조

```
block-chain/
├── main.go               # 애플리케이션 실행 진입점
├── environment.toml      # 설정파일
├── app/                  # 앱 실행 흐름 제어
├── config/               # TOML 기반 환경 로딩
├── global/               # 전역 변수 및 상태 관리
├── repository/           # 데이터 정의 (Block, Wallet 등)
├── service/              # 핵심 로직 (POW, 트리, 지갑 등)
├── types/                # 타입 및 상수 정의
└── .idea/                # 개발 IDE 설정
```

---

##  핵심 기능 및 주요 코드

###  1. 블록 생성 및 체인 구성

```go
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}
```

```go
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
```

- 블록 간 해시는 SHA-256 기반으로 연결됩니다.
- 생성 시 POW를 적용하여 유효한 해시값을 찾습니다.

---

### 2. 작업 증명 (Proof of Work)

```go
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}
	return nonce, hash[:]
}
```

- 목표 난이도보다 작은 해시가 나올 때까지 반복 계산
- 블록의 무결성과 생성 난이도를 보장하는 핵심 알고리즘입니다.

---

###  3. 머클 트리 구성

```go
func NewMerkleTree(data [][]byte) *MerkleTree {
	// 트랜잭션 데이터를 해시하여 트리 구성
}
```

- 블록 내 트랜잭션들을 머클 트리로 묶어 루트 해시를 생성
- 전체 블록을 재검증하지 않고도 특정 트랜잭션의 무결성을 확인할 수 있습니다.

---

### 4. 지갑 및 키 생성

```go
func NewWallet() *Wallet {
	private, public := newKeyPair()
	address := generateAddress(public)
	return &Wallet{PrivateKey: private, PublicKey: public, Address: address}
}
```

- 개인키/공개키는 ECDSA 기반으로 생성
- 공개키는 해시 처리를 통해 지갑 주소로 변환됩니다

```go
wallet := NewWallet()
fmt.Println("Your address:", wallet.Address)
```

---

## 실행 방법

```bash
# 환경 파일 수정 (필요시)
vim environment.toml

# 의존성 정리
go mod tidy

# 프로젝트 실행
go run main.go
```

---

## 학습 포인트

- SHA-256 기반 블록 연결 및 해시 검증 로직 직접 구현
- POW 알고리즘의 작동 방식과 보안성을 코드로 체득
- 머클 트리를 통한 트랜잭션 무결성 보장 방식 이해
- ECDSA 서명 및 주소 생성 로직을 통한 블록체인 지갑 구조 학습
- Go 언어의 모듈화 구조 및 TOML 기반 환경 설정 실습

---

## 향후 확장 계획

- [ ] 트랜잭션 풀 및 서명 검증 로직 구현
- [ ] REST API 또는 CLI 인터페이스 개발
- [ ] 블록 탐색기 (Block Explorer) 웹 대시보드 구축
- [ ] 스마트 컨트랙트 기반 실행 엔진 추가
- [ ] P2P 멀티 노드 블록체인 네트워크 구성

---

## 기술적인 후기

- SHA-256 기반 블록 해시 및 머클 트리 구성 과정을 직접 구현하며 데이터 무결성과 블록 연결 방식에 대한 명확한 이해를 얻었습니다.
- POW 알고리즘 구현을 통해 블록체인의 합의 메커니즘과 연산 비용에 대해 깊이 체감할 수 있었습니다.
- ECDSA 키 생성과 주소 변환 과정을 직접 설계해보며 실제 지갑 구조를 실습했습니다.
- 프로젝트 구조 설계 및 TOML 설정을 도입하며, 서비스 확장성과 유지보수성을 고려한 설계 역량을 키웠습니다.

---

## 참고 자료

- [Mastering Bitcoin - Andreas M. Antonopoulos](https://github.com/bitcoinbook/bitcoinbook)
- Go 공식 문서 및 RFC


---