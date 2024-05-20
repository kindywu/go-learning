package main

// 导入所需的库
import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// 可以设置一些参数，比如节点ID
var (
	nodeID    = flag.Int("id", 0, "node ID")
	electName = flag.String("name", "my-test-elect", "election name")
)

var count int

// 选主
func elect(e1 *concurrency.Election, electName string) {
	log.Println("acampaigning for ID:", *nodeID, electName)
	// 调用Campaign方法选主,主的值为value-<主节点ID>-<count>
	if err := e1.Campaign(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("campaigned for ID:", *nodeID)
	count++
}

// 为主设置新值
func proclaim(e1 *concurrency.Election, electName string) {
	log.Println("proclaiming for ID:", *nodeID, electName)
	// 调用Proclaim方法设置新值,新值为value-<主节点ID>-<count>
	if err := e1.Proclaim(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("proclaimed for ID:", *nodeID)
	count++
}

// 重新选主，有可能另外一个节点被选为了主
func resign(e1 *concurrency.Election, electName string) {
	log.Println("resigning for ID:", *nodeID, electName)
	// 调用Resign重新选主
	if err := e1.Resign(context.TODO()); err != nil {
		log.Println(err)
	}
	log.Println("resigned for ID:", *nodeID)
}

// 查询主的信息
func query(e1 *concurrency.Election, electName string) {
	// 调用Leader返回主的信息，包括key和value等信息
	resp, err := e1.Leader(context.Background())
	if err != nil {
		log.Printf("failed to get the current leader: %v electName: %s", err, electName)
	} else {
		log.Println("current leader:", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value), electName)
	}
}

// 可以直接查询主的rev信息
func rev(e1 *concurrency.Election, electName string) {
	rev := e1.Rev()
	log.Println("current rev:", rev, electName)
}

func watch(e1 *concurrency.Election, electName string) {
	ch := e1.Observe(context.TODO())

	log.Println("start to watch for ID:", *nodeID, electName)
	for i := 0; i < 10; i++ {
		resp := <-ch
		log.Println("leader changed to", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
	}
}

func main() {
	flag.Parse()

	// 将etcd的地址解析成slice of string
	cli, err := clientv3.New(clientv3.Config{
		// Endpoints: []string{"etcd:2379"},
		Endpoints:   []string{"etcd-etcd1-1:2379", "etcd-etcd2-1:22379", "etcd-etcd3-1:32379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 创建session,如果程序宕机导致session断掉，etcd能检测到
	session, _ := concurrency.NewSession(cli)
	defer session.Close()

	// 生成一个选举对象。下面主要使用它进行选举和查询等操作
	// 另一个方法ResumeElection可以使用既有的leader初始化Election
	e1 := concurrency.NewElection(session, *electName)

	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		switch action {
		case "elect": // 选举命令
			go elect(e1, *electName)
		case "proclaim": // 只更新leader的value
			proclaim(e1, *electName)
		case "resign": // 辞去leader,重新选举
			resign(e1, *electName)
		case "watch": // 监控leader的变动
			go watch(e1, *electName)
		case "query": // 查询当前的leader
			query(e1, *electName)
		case "rev":
			rev(e1, *electName)
		default:
			fmt.Println("unknown action")
		}
	}
}
