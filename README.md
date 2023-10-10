Quant1x-Engine
===

量化交易（数据类）引擎

## 1. 设计原则
- 更新历史数据，盘后更新历史数据
- 回补历史数据，对于新增特征组合或者因子要能回补历史数据，以便回测
- 盘中更新数据，盘中决策很重要，特征组合要有根据5档行情即时数据进行增量计算的能力
- 缓存数据必须具备按照日期切换的功能

## 2. 模块划分

| 级别 | 模块       | 功能   | 盘中更新数据 | 更新当日数据 | 回补历史数据 |
|:---|:---------|:-----|:-------|:-------|:-------|
| 0  | cache    | 数据缓存 | [ ]    | [ ]    | [ ]    |
| 0  | factors  | 量化因子 | [ ]    | [ ]    | [ ]    |
| 0  | features | 特征   | [ ]    | [ ]    | [ ]    |
| 0  | tracker  | 回测   | [ ]    | [ ]    | [ ]    |

## 3. 使用示例

### 3.1 更新数据

```shell
engine update --all
```

### 3.2 补登历史特征数据

```shell
engine repair --history --start=20230101
```

### 3.3 执行1号策略

```shell
engine --strategy=1
```

## 4. 协同开发

### 4.1 git仓库规则

#### 4.1.1 分支规则

主库的分支分 ***master*** 和 ***子版本***

| 级别 | 分类                | 说明     | 示例                  | 
|:---|:------------------|:-------|:--------------------|
| 0  | master            | 主线版本分支 |                     |
| 1  | {Major}.{Minor}.x | 子版本分支  | 主版本1,子版本2, 分支为1.2.x |

#### 4.1.2 tag规则

- tag只有一条, 小写字母v开头, 然后依次是主版本号、子版本号和修订版本号, 中间用“.”分隔。
- 格式:

```shell
v${Major}.${Minor}.${Revision}
```

### 4.2 约定

| 级别 | 分类  | 单项           | 说明                           | 示例                                                   | 
|:---|:----|:-------------|:-----------------------------|:-----------------------------------------------------|
| 0  | 关键词 | securityCode | 完整的证券代码,格式:{MARKET_ID}{CODE} | 上证指数: sh000001</br>永鼎股份: sh600105</br>鼎汉技术: sz300011 |

### 4.3 协同开发流程

- fork项目到自己的git仓库
- 先提issue, 再实现、测试
- 提交PR
- PR审核, 不通过, 驳回
- PR审核, 通过, 按照不同的issue级别, 合并当前次级版本分支, 并发布版本