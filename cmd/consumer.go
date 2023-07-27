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
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/spf13/cobra"
)

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "作为consumer消费主题数据",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		nameServer, _ := cmd.Flags().GetString("name")
		group, _ := cmd.Flags().GetString("group")
		topic, _ := cmd.Flags().GetString("topic")
		fmt.Printf("启动订阅服务 NameServer: %s,Topic: %s, GroupName: %s\n", nameServer, topic, group)
		s := NewSubcriber([]string{nameServer}, topic, group)
		s.Start()
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)
	var nameServer string
	consumerCmd.Flags().StringVarP(&nameServer, "name", "n", "", " name server addr")
	var groupName string
	consumerCmd.Flags().StringVar(&groupName, "group", "tester", "consumer group name")
	var TopicName string
	consumerCmd.Flags().StringVar(&TopicName, "topic", "test-topic", "topic name")
	consumerCmd.MarkFlagRequired("name")

}

type Subcriber struct {
	client rocketmq.PushConsumer
}

func NewSubcriber(server []string, topicName, grpName string) Subcriber {

	// 订阅主题、消费
	// 创建一个consumer实例
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer(server),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(grpName),
	)
	if err != nil {
		fmt.Println("创建consumer实例错误:", err)
	}
	// 订阅topic
	err = c.Subscribe(topicName, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("接收到消息: %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		fmt.Printf("订阅消息错误: %s\n", err.Error())
	}
	s := Subcriber{}
	s.client = c
	return s
}

func (s *Subcriber) Start() {
	// 启动consumer
	err := s.client.Start()

	if err != nil {
		fmt.Printf("启动Consumer错误: %s\n", err.Error())
		os.Exit(-1)
	}
	s.hookSignals()
	for {
	}
}

func (s *Subcriber) Stop() {
	err := s.client.Shutdown()
	if err != nil {
		fmt.Printf("停止Consumer错误: %s\n", err.Error())
	}
}

func (s *Subcriber) hookSignals() {
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
