# Yu-Gi-Oh! AI Tools (游戏王 AI 分析工具集)

[English](README.md) | **简体中文**

基于 AI Agent Skills 的游戏王卡组自动化分析工具。

## 项目简介

本项目提供一系列工具和 Agent Skills，旨在辅助 Gemini CLI 或 Claude Code 等 AI 助手深度分析游戏王卡组。通过自动化解析 `.ydk` 文件、计算起手概率并生成多维度分析报告，为 AI 赋予专业级的卡组评估能力。

## 项目结构

```text
yu-gi-oh-ai-tools/
├── skills/                    # 适配 Gemini CLI / Claude Code 的 Agent Skills
│   ├── ygo-deck-analyze/      # 核心编排逻辑，负责协调整个分析流程
│   ├── deck-parser/           # 卡组分类与概率计算引擎
│   ├── deck-analyzer/         # 强度分析与对局评估
│   ├── starter-probability/   # 超几何分布概率计算器
│   ├── count-deck/            # 卡组张数统计工具
│   └── simple-calculator/     # 基础数学工具
│
├── ygo-db/                    # 基于 Go 的卡牌数据库 CLI
│   └── cmd/ygo-db-cli/        # 高性能卡片查询工具
│
└── test-dir-zh-cn/            # 测试目录，包含示例卡组
```

## 功能特性

### 核心分析 Skills

| Skill | 描述 |
|-------|-------------|
| **ygo-deck-analyze** | 分析流程总控。协调从原始数据获取到最终报告生成的完整流水线。 |
| **deck-parser** | 事实提取引擎。识别系统轴、展开点、补点及手坑，并计算基础展开概率。 |
| **deck-analyzer** | 强度评估引擎。分析展开多样性、终场强度、卡点 (Chokepoint) 及资源循环。 |

### 实用工具 Skills

| Skill | 描述 |
|-------|-------------|
| **starter-probability** | 使用超几何分布计算在 5 张起手手中抽到特定组合的概率。 |
| **count-deck** | 快速统计 YDK 文件中主卡组、额外卡组和副卡组的张数。 |
| **ygoprodeck-get-meta** | 从 YGOPRODeck API 获取最新环境的主流卡组。 |

### 数据库工具

- **ygo-db-cli**: 基于 Go 语言编写的高性能命令行工具，支持按 ID、名称、系列或 YDK 文件查询卡片。

## 分析流水线

`ygo-deck-analyze` Skill 将自动执行以下流程：

1.  **数据采集**: 通过 `ygo-db-cli` 获取原始卡片元数据。
2.  **结构解析**: 识别展开点、核心组件和防御性卡片。
3.  **强度评估**: 对卡组稳定性、上限等多个维度进行量化评分。
4.  **生成报告**: 汇总结果并生成格式化的 `report.md`。

## 输出成果

每次分析将生成以下文件：
- `deck_raw.csv`: 从数据库导出的原始元数据。
- `deck_parsed.json`: 结构化卡组信息（系统轴、展开点、手坑分布）。
- `deck_analysis.json`: 多维度评分与详细分析结论。
- `report.md`: 供人类阅读的分析报告及战术建议。

## 六维评分体系

| 维度 | 评估重点 |
|-----------|-------------|
| **稳定性 (Consistency)** | 起手拿到展开点的概率。 |
| **上限 (Ceiling)** | 无干扰情况下的终端压制力。 |
| **抗干扰 (Resilience)** | 吃手坑后的补点能力与妥协场。 |
| **后攻 (Going Second)** | 突破对方场面及返场能力。 |
| **续航 (Grind)** | 资源回收与长线作战能力。 |
| **多样性 (Diversity)** | 展开路线的灵活性与多变性。 |

## 快速上手

```bash
# 生成卡组的原始数据
ygo-db-cli get-cards-by-ydk -i deck.ydk -o deck_raw.csv

# 在 Gemini CLI / Claude Code 中运行
# 激活 'ygo-deck-analyze' Skill 处理上述生成的 CSV
```

## 环境要求

- **AI Skills**: Gemini CLI、Claude Code 或兼容的 AI 助手环境。
- **ygo-db**: Go 1.21+
- **Python**: 3.14+ (部分工具脚本需要)

## License

Apache 2.0

---

*游戏王是科乐美数字娱乐公司的商标。本项目与科乐美数字娱乐公司无关联、未被其认可、赞助或特别批准。*
