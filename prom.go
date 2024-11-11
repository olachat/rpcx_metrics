package prom

//
// // Start prometheus服务启动
// func Start(ctx context.Context, serviceName string, consulAddr []string, wg *sync.WaitGroup) {
// 	defer wg.Done()
//
// 	availablePort := tool.GetAvailablePort()
// 	log.Printf("serviceName: %+v prometheus port: %+v start", serviceName, availablePort)
//
// 	http.Handle("/metric", DefaultRegisterer().GetMetricAPIHandler())
// 	// 增加健康检测
// 	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
// 		writer.WriteHeader(http.StatusOK)
// 		_, _ = writer.Write([]byte("pong"))
// 	})
//
// 	srv := &http.Server{
// 		Addr: fmt.Sprintf(":%d", availablePort),
// 	}
// 	consulDiscovery := NewConsulDiscovery("go_rpc_exporter", availablePort, consulAddr)
//
// 	go func() {
// 		// 向consul进行注册
// 		if err := consulDiscovery.Register([]string{serviceName}, map[string]string{
// 			"service_name": serviceName,
// 		}); err != nil {
// 			panic(err)
// 		}
// 		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			panic(fmt.Sprintf("prometheus ListenAndServer failed: %v", err))
// 		}
// 	}()
//
// 	<-ctx.Done()
// 	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
//
// 	// 解除服务
// 	_ = consulDiscovery.Deregister()
//
// 	if err := srv.Shutdown(shutdownCtx); err != nil {
// 		log.Printf("prometheus Shutdown failed err: %+v", err)
// 	}
// 	log.Printf("prometheus shutdown finished")
// }
