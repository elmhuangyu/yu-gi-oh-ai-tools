---
name: ygo-deck-analyze
description: 游戏王卡组完整分析流程的 parent skill。当用户提供 .ydk 文件并要求分析卡组时使用此 skill。负责建立工作目录、追踪进度、按顺序调用 deck-parser 和 deck-analyzer，并生成最终人类可读报告。触发词：分析卡组、.ydk、卡组分析、帮我看看这套卡。
---

# ygo-deck-analyze

## 职责

这是整个分析流程的协调者（orchestrator）。

负责：
- 建立工作目录和进度追踪文件
- 按顺序执行四个步骤
- 每步完成后将产物 dump 成文件
- 如果某步失败，记录断点，下次可从断点恢复
- 生成最终人类可读报告

不负责：
- 卡片信息获取（由 deck-parser 处理）
- 对局分析判断（由 deck-analyzer 处理）

---

## 工作目录结构

每个卡组建立独立目录，命名规则为 `{ydk文件名}_analysis/`。

```
{deck_name}_analysis/
├── TODO.md              # 进度追踪，每步完成后更新
├── deck_raw.csv         # Step 1 产物：CLI 原始返回数据（CSV格式）
├── deck_parsed.json     # Step 2 产物：deck-parser 输出
├── deck_analysis.json   # Step 3 产物：deck-analyzer 输出
└── report.md            # Step 4 产物：最终人类可读报告
```

---

## 断点恢复逻辑

开始前，先检查工作目录是否已存在：

- 如果目录不存在 → 全新开始，执行 Step 1
- 如果目录存在，读取 `TODO.md` → 找到第一个未完成的步骤，从该步骤继续
- 如果某个产物文件存在且非空 → 该步骤视为已完成，跳过

---

## 执行流程

### Step 1 — 初始化工作目录

1. 根据 `.ydk` 文件名创建工作目录

```bash
mkdir -p {deck_name}_analysis
```

2. 创建 `TODO.md`，内容如下：

```markdown
# {deck_name} 分析进度

## 文件
- ydk: {ydk_filename}
- 开始时间: {timestamp}

## 步骤
- [ ] Step 1: 获取原始卡片数据 → deck_raw.csv
- [ ] Step 2: 解析卡组结构 → deck_parsed.json
- [ ] Step 3: 对局分析 → deck_analysis.json
- [ ] Step 4: 生成报告 → report.md
```

3. 在 `TODO.md` 中将 Step 1 标记为进行中。

---

### Step 2 — 获取原始卡片数据

使用 `ygo-db-cli` CLI 工具获取卡片信息：

```bash
ygo-db-cli get-cards-by-ydk -i {ydk文件路径} -o deck_raw.csv
```

CLI 工具会将数据（CSV格式）完整 dump 到 `deck_raw.csv`。

> 保存原始数据的目的：如果后续解析逻辑有问题，可以直接重跑解析步骤，不需要重新调用 CLI。

完成后更新 `TODO.md`：
```markdown
- [x] Step 1: 获取原始卡片数据 → deck_raw.csv
```

---

### Step 3 — 解析卡组结构

读取 `deck_raw.csv`，调用 `deck-parser` skill 进行卡组结构解析。

输入：`deck_raw.csv`
输出：按 `deck-parser` skill 定义的 JSON schema 生成结果，dump 到 `deck_parsed.json`。

完成后更新 `TODO.md`：
```markdown
- [x] Step 2: 解析卡组结构 → deck_parsed.json
```

---

### Step 4 — 对局分析

读取 `deck_raw.csv` 和 `deck_parsed.json`，调用 `deck-analyzer` skill 进行完整对局分析。

输入：`deck_raw.csv` + `deck_parsed.json`
输出：按 `deck-analyzer` skill 定义的 JSON schema 生成结果，dump 到 `deck_analysis.json`。

完成后更新 `TODO.md`：
```markdown
- [x] Step 3: 对局分析 → deck_analysis.json
```

---

### Step 5 — 生成人类可读报告

读取 `deck_parsed.json` 和 `deck_analysis.json`，生成自然语言分析报告，dump 到 `report.md`。

报告格式见下方「报告模板」。

完成后更新 `TODO.md`：
```markdown
- [x] Step 4: 生成报告 → report.md

## 完成
- 完成时间: {timestamp}
```

---

## 报告模板

`report.md` 必须使用以下结构：

```markdown
# {archetype} 卡组分析报告

**卡组类型**：{deck_type} ｜ **操作难度**：{difficulty}  
**主卡组**：{main}张 ｜ **额外卡组**：{extra}张 ｜ **副卡组**：{side}张

---

## 卡组概览

{用2-3句话描述卡组的核心机制和风格}

**主引擎**：{primary_engine}  
**副引擎 / Package**：{secondary_engines}

---

## 起手稳定性

- Starter 开出率：{starter_open_rate}%
- 手坑开出率：{handtrap_open_rate}%

{用1-2句话说明稳定性的含义，例如"平均每4手有3手能正常启动"}

---

## 典型先攻终场

**终场类型**：{board_type}  
**妨碍数量**：{interaction_count}

| 终场卡片 | 作用 |
|---|---|
{终场卡列表}

**抗性**：{抗性说明}  
**自锁风险**：{self_lock_risk}

> 如果终场类型为 `RESOURCE_DENIAL_BOARD`，在此说明对手的哪些资源会被破坏，以及通常在第几回合让对手断粮。

---

## 展开多样性

**整体评价**：{overall_verdict} — {overall_reason}

| 干扰类型 | 应对能力 | 说明 |
|---|---|---|
{干扰类型列表}

---

## 终场弱点

**最大威胁**：{biggest_threat} — {biggest_threat_reason}

**盘面突破手段**

| 去除手段 | 有效程度 | 原因 |
|---|---|---|
{盘面突破列表}

**机制封锁手段**

| 封锁类型 | 有效程度 | 原因 |
|---|---|---|
{机制封锁列表}

---

## 系统性弱点

{对每个系统性弱点用一段话说明：依赖什么机制、哪些卡能针对、被针对后还能做什么、这类卡在当前 meta 中是否常见}

---

## 手坑敏感度

| 手坑 | 影响程度 | 被打断时机 |
|---|---|---|
{手坑列表}

---

## 关键 Chokepoint

{对每个 chokepoint 用一段话说明：是什么、为什么关键、怕什么、有没有替代路}

---

## 后攻能力

{用自然语言描述后攻突破手段、效率、突破后余力、OTK 潜力}

---

## 资源续航

{用自然语言描述检索机制、墓地循环、面对增殖G的策略、长局能力}

---

## 六维评分

| 维度 | 分数 | 说明 |
|---|---|---|
| 起手稳定性 consistency | {value}/10 | {reason} |
| 终场上限 ceiling | {value}/10 | {reason} |
| 抗干扰能力 resilience | {value}/10 | {reason} |
| 后攻能力 going_second | {value}/10 | {reason} |
| 资源续航 grind | {value}/10 | {reason} |
| 展开多样性 route_diversity | {value}/10 | {reason} |

---

## 总结

**优势**
{优势列表}

**弱点**
{弱点列表，包含具体卡名和场景}
```

---

## 错误处理

如果某个步骤失败：
1. 停止执行，告知用户哪步失败、失败原因、以及已保存的产物文件路径
2. 用户下次重新触发时，从失败步骤重新开始
