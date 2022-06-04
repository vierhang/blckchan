package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

// PrintUsage 用法展示
func PrintUsage() {
	fmt.Println("Usage:")
	//初始化区块链
	fmt.Printf("\tcreageblockchain -address address -- 创建区块链 \n")
	// 添加区块
	fmt.Printf("\taddblock -data DATA --添加区块\n")
	// 打印完整区块信息
	fmt.Printf("\tprintchain --输出区块链信息\n")
}

//初始化区块链
func (c *CLI) CreateBlockChain(address string) {
	CreateBlockChainWithGenesisBLock(address)
}

//添加区块
func (c *CLI) AddBlock(txs []*Transaction) {
	if !dbExists() {
		fmt.Println("dbExists")
		os.Exit(1)
	}
	NewBlockChain().AddBlock(txs)
}

//打印完整区块链信息
func (c *CLI) PrintChain() {
	if !dbExists() {
		fmt.Println("dbExists")
		os.Exit(1)
	}
	NewBlockChain().PrintChan()
}

// IsValidArgs 参数数量检测函数
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		fmt.Println("os.Args < 2")
		os.Exit(1)
	}
}

// Run 命令行运行函数
func (c *CLI) Run() {
	IsValidArgs()
	// 新建相关命令
	// 添加区块
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBLCWithGenesisBlockCmd := flag.NewFlagSet("cerateblockchain", flag.ExitOnError)

	// 数据参数处理
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块数据")
	// 创建区块时指定的矿工地址
	flagCreageBlockchainArg := createBLCWithGenesisBlockCmd.String("address", "weihang", "指定接受系统奖励的矿工地址")
	// 判断命令
	switch os.Args[1] {
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse addBlockCmd err %v \n", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse printChainCmd err %v \n", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse createBLCWithGenesisBlockCmd err %v \n", err)
		}
	default:
		PrintUsage()
		fmt.Println(os.Args[1])
		fmt.Println("os.Args[1] switch error")
		os.Exit(1)
	}
	// 添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		c.AddBlock([]*Transaction{})
	}
	// 输出区块信息
	if printChainCmd.Parsed() {
		c.PrintChain()
	}
	// 创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed() {
		if *flagCreageBlockchainArg == "" {
			fmt.Println("flagCreageBlockchainArg == ''")
			PrintUsage()
			os.Exit(1)
		}
		c.CreateBlockChain(*flagCreageBlockchainArg)
	}
}
