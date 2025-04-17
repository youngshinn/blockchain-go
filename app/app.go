package app

import (
	"block-test/config"
	"block-test/global"
	"block-test/repository"
	"block-test/service"
	"block-test/types"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hacpy/go-ethereum/common"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/inconshreveable/log15"
)

type App struct {
	config     *config.Config
	service    *service.Service
	repository *repository.Repository

	log log15.Logger
}

func NewApp(config *config.Config) {
	a := &App{
		config: config,
		log:    log15.New("module", "app"),
	}

	var err error

	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	} else {
		a.service = service.NewService(a.repository, 1)

		a.log.Info("Module start", "time", time.Now().Unix())

		sc := bufio.NewScanner(os.Stdin)

		useCase()

		for {
			from := global.FROM()

			if from != "" {
				a.log.Info("Current Connected wallet", "from", from)
				fmt.Println()
			}

			sc.Scan()
			fmt.Println(sc.Text())
			input := strings.Split(sc.Text(), " ")
			if err = a.inputValue(input); err != nil {
				a.log.Error("Failed to cli", "err", err, "input", input)
				fmt.Println()
			}
		}
	}
}

// 들어오는 값을 분기처리
func (a *App) inputValue(input []string) error {
	if len(input) == 0 {
		msg := errors.New("check use case")
		panic(msg)
	} else {

		from := global.FROM()

		switch input[0] {
		case types.CreateWallet:

			if wallet := a.service.MakeWallet(); wallet == nil {
				panic("failed to create wallet")
			} else {
				fmt.Println()
				fmt.Println("success to create wallet")
				a.log.Info("CreatedWallett", "pk:", wallet.PrivateKey, "pu", wallet.PublicKey)
				fmt.Println("")
			}

		case types.TransferCoin:

			if from == "" {
				fmt.Println()
				a.log.Debug("Request value, to is unCoreect")
				fmt.Println()

			} else if input[1] == "" || input[2] == "" {
				fmt.Println()
				a.log.Debug("Not Connected Wallet Connect Wallet First")
				fmt.Println()
			} else {
				// from, to, value
				a.service.CreateBlock(from, input[1], input[2])
			}

		case types.MintCoin:

			if input[1] == "" || input[2] == "" {
				fmt.Println()
				a.log.Debug("Request value, to is unCoreect")
				fmt.Println()

			} else {
				// common.Address
				a.service.CreateBlock((common.Address{}).String(), input[1], input[2])
			}

			if from == "" {
				a.log.Debug("Connect Wallet First")
				fmt.Println()
			} else {
				a.service.CreateBlock(from, input[1], input[2])
			}

		case types.ConnectWallet:
			fmt.Println()

			if from != "" {
				a.log.Debug("Already have wallet", "from", from)
				fmt.Println()
			} else {
				if wallet, err := a.service.GetWallet(input[1]); err != nil {
					if err == mongo.ErrNoDocuments {
						a.log.Debug("failed to get wallet pk is nil", "pk", input[1])
					} else {
						a.log.Crit("failed to Find Wallet", "pk", input[1], "err", err)
					}
				} else {
					global.SetFrom(wallet.PublicKey)
					fmt.Println()
					a.log.Info("success to get wallet", "pk", wallet.PublicKey)
					fmt.Println()
				}
			}

		case types.ChangeWallet: //기존에 설정된 FROM 값(즉 DB에 저당되어잇는 값만 호출)
			fmt.Println()

			if from == "" {
				a.log.Debug("Connect Wallet First")
				fmt.Println()
			} else {
				if wallet, err := a.service.GetWallet(input[1]); err != nil {
					if err == mongo.ErrNoDocuments {
						a.log.Debug("failed to get wallet pk is nil", "pk", input[1])
					} else {
						a.log.Crit("failed to Find Wallet", "pk", input[1], "err", err)
					}
				} else {
					global.SetFrom(wallet.PublicKey)
					fmt.Println()
					a.log.Info("success to change Wallet", "pk", wallet.PublicKey)
					fmt.Println()
				}
			}
		default:
			return errors.New("Unknown input")
		}

	}

	return nil
}

func useCase() {
	fmt.Println()

	fmt.Println("This is module for block-chain mort mongo")
	fmt.Println()
	fmt.Println("Use Case")
	fmt.Println("1. ", types.CreateWallet)
	fmt.Println()
	fmt.Println("2. ", types.ConnectWallet, "<PK>")
	fmt.Println()
	fmt.Println("3. ", types.ChangeWallet, "<PK>")
	fmt.Println()
	fmt.Println("4. ", types.TransferCoin, " <to>")
	fmt.Println()
	fmt.Println("5. ", types.MintCoin, "<to>")
} //사용자 명령어 선택
