# layeraoi(分层AOI)
layeraoi是实现的一个AOI算法容器（目前最大容量64层），设计的思想是依据单位的不同层级的视野大小，看到不同的范围的对象。例如刺激战场的主角可以假装分为三层视野：
第一层是能够看到的武器视野
第二层是能够看到的空投补给箱范围
第三层是敌方玩家范围
经过分层后，主角只需要获得他视野内应该关注的对象，从而减少对象信息的收发。


layeraoi除了实现分层类型，还实现了单层内的动态分层。即单个层内，触发了负载的条件，那么这个层是会被分割的，即一个层变成两个层，同样的，一个单类层内的对象数目
太少也可能会触发层合并。
演示地址：http://121.5.223.223/wasm_exec.html

简单的图形
![image](https://github.com/hitong/layeraoi/blob/main/awesome/base.png)

不简单的图形
![image](https://github.com/hitong/layeraoi/blob/main/awesome/base2.png)
