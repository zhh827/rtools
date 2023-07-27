/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/spf13/cobra"
)

// producerCmd represents the producer command
var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "作为producer向topic发送数据",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		nameServer, _ := cmd.Flags().GetString("name")
		group, _ := cmd.Flags().GetString("group")
		topic, _ := cmd.Flags().GetString("topic")
		fmt.Printf("启动消息推送服务 NameServer: %s,Topic: %s, GroupName: %s\n", nameServer, topic, group)
		s := NewSender([]string{nameServer}, topic, group)
		s.Start()
	},
}

func init() {
	rootCmd.AddCommand(producerCmd)
	var nameServer string
	producerCmd.Flags().StringVarP(&nameServer, "name", "n", "", " name server addr")
	var groupName string
	producerCmd.Flags().StringVar(&groupName, "group", "tester", "consumer group name")
	var TopicName string
	producerCmd.Flags().StringVar(&TopicName, "topic", "test-topic", "topic name")
	producerCmd.MarkFlagRequired("name")

}

type Sender struct {
	client    rocketmq.Producer
	topicName string
}

func NewSender(server []string, topicName, grpName string) Sender {
	// 发送消息
	// 创建一个producer实例
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(server),
		producer.WithRetry(2),
		producer.WithGroupName(grpName),
	)
	if err != nil {
		fmt.Println("创建 producer错误:", err)
	}
	s := Sender{}
	s.client = p
	s.topicName = topicName
	return s
}

func (s *Sender) Start() {
	s.hookSignals()
	// 启动
	err := s.client.Start()
	if err != nil {
		fmt.Printf("启动producer错误: %s", err.Error())
		os.Exit(1)
	}
	var msg string
	for {
		fmt.Scanln(&msg) //
		// 发送消息
		result, err := s.client.SendSync(context.Background(), &primitive.Message{
			Topic: s.topicName,
			Body:  []byte(msg),
		})
		if err != nil {
			fmt.Printf("发送消息错误: %s\n", err.Error())
		} else {
			fmt.Printf("发送消息成功: result=%s\n", result.String())
		}
	}

}

func (s *Sender) Stop() {
	err := s.client.Shutdown()
	if err != nil {
		fmt.Printf("停止Consumer错误: %s\n", err.Error())
	}
}

func (s *Sender) hookSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			switch sig {
			case syscall.SIGTERM:
				fallthrough
			case syscall.SIGINT:
				fmt.Printf("received signal %s, exiting...", sig.String())
				s.Stop()
				os.Exit(0)
			}
		}
	}()
}
