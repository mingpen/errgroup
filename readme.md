# errgroup

基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

# note

1. errgroup给我们提供了go routine生命周期的管理工具；
2. `g, ctx := errgroup.WithContext(context.Background())`  
里面的ctx是用来监听是否有一个go routine已经出错，并退出了；
3. 要响应group内的routine异常，免不了 `select { case <-ctx.Done():}`,并退出；
4. 对于封装好的阻塞同步函数【例如http ListenAndServe】，
可以另外启动一个go routine监听group内的异常，并调用对应package的Shutdown、Stop函数；
5. http\signal 官方库都使用不同的 routine生命周期管理；
6. http使用map 存储conn的状态，shoudown时检查conn状态，来判断是否关闭完；
7. signal是通过调用stop方法来管理 自己package创建的
