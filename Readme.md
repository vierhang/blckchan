# viehrang-仿公链

# day1
1. 生成区块基本模块
2. 新建区块
3. 如何生成哈希
4. 类型转换

# day2
1. 实现链表（简单版用切片实现）
2. 实现上链

# da3 实现共识Pow
1. pow 结构分析
2. 设置目标难度值
3. 哈希碰撞
4. 数据准备

# day4 持久化
1. https://github.com/boltdb/bolt
2. 定义数据库名和表名
3. biltdb基本操作（insert、read）
4. 遍历区块链
5. 迭代器

# day5 实现命令行

# day6 
1. 获取blockchain 对象

# day 7 修改交易结构
1. 实现交易结构替换Data
2. 实现输入、输出结构

# day8 实现coinBase（挖矿-没有输入）
1. 实现创世块挖矿
2. 交易哈希序列化的实现


# day9 实现cli发起转账
1. 添加send命令行
2. 多笔交易json  JSONToArray
3. `bcli send -from "[\"test\"]" -to "[\"b\"]" -amount "[\"20\"]"`


# day10 实现通过挖矿生成新区块
1. 实现挖矿功能
2. 通过交易调用挖矿

# day 11 实现普通交易
