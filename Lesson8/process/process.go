package main

import (
	_ "gocloud.dev/pubsub/kafkapubsub"

	"context"
	"log"
	"math/rand"
	"time"

	"github.com/mediocregopher/radix/v3"
	"gocloud.dev/pubsub"
)

var (
	connFunc = func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(time.Second),
		)
	}
	s   *radix.Pool
	sub *pubsub.Subscription
)

func storage() *radix.Pool {
	if s != nil {
		return s
	}
	var err error
	s, err = radix.NewPool("tcp", "redis:6379", 1, radix.PoolConnFunc(connFunc))
	if err != nil {
		panic(err)
	}
	return s
}

func subscription() (*pubsub.Subscription, error) {
	if sub != nil {
		return sub, nil
	}
	var err error
	sub, err = pubsub.OpenSubscription(context.Background(), "kafka://process?topic=rates")
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func main() {
	for {
		time.Sleep(10 * time.Second)
		s, err := subscription()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		msg, err := s.Receive(context.Background())
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}

		err = storage().Do(radix.Cmd(nil, "LPUSH", "result", string(msg.Body)))
		if err != nil {
			log.Println(err)
		}
		if rand.Float64() < .05 {
			_ = storage().Do(radix.Cmd(nil, "LTRIM", "result", "0", "9"))
		}
		msg.Ack()
	}
}
