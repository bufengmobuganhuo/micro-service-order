package main

import (
	common "github.com/bufengmobuganhuo/micro-service-common"
	"github.com/bufengmobuganhuo/order/domain/repository"
	service2 "github.com/bufengmobuganhuo/order/domain/service"
	"github.com/bufengmobuganhuo/order/handler"
	order "github.com/bufengmobuganhuo/order/proto/order"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

var (
	QPS = 1000
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/config")
	if err != nil {
		log.Fatal(err)
	}

	// 注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪，jaeger
	t, io, err := common.NewTracer("go.micro.service.order", "localhost:8631")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 初始化数据库
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	db, err := mysqlInfo.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SingularTable(false)

	//repo := repository.NewOrderRepository(db)
	//repo.InitTable()

	// 创建实例
	orderDataService := service2.NewOrderDataService(repository.NewOrderRepository(db))

	// prometheus, 暴露监控地址
	common.PrometheusBoot(9092)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.order"),
		micro.Version("latest"),
		// 暴露的服务地址
		micro.Address(":9085"),
		// 注册中心
		micro.Registry(consul),
		// 链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
		// 添加监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// Initialise service
	service.Init()

	// Register Handler
	order.RegisterOrderHandler(service.Server(), handler.Order{OrderDataService: orderDataService})

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
