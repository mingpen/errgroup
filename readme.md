# errgroup

基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

# note

1. errgroup给我们提供了go routine生命周期的管理工具；
2. `g, ctx := errgroup.WithContext(context.Background())`  
里面的ctx是用来监听是否有一个go routine已经出错，并退出了；
3. 要响应group内的routine异常，免不了 `select { case <-ctx.Done():}`；
4. http\signal 官方库都使用不同的 routine生命周期管理；
5. http使用map 存储conn的状态，shoudown时检查conn状态，来判断是否关闭完；
6. signal是通过调用stop方法来管理 自己package创建的
