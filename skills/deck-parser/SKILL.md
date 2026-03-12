---
name: deck-parser
description: 游戏王卡组结构解析 skill。接收 MCP 原始卡片数据（deck_raw.json），识别 engine、starter、extender、手坑，计算起手概率，输出结构化卡组数据 JSON。只做事实提取，不做对局判断。通常由 yugioh-deck-analyze 调用，输入为 deck_raw.json，输出为 deck_parsed.json。
---

# deck-parser

## 职责边界

这个 skill **只做事实提取**，不做判断。

- ✅ 卡片是什么、效果是什么、分几张
- ✅ 这张卡属于哪个 engine
- ✅ 这张卡是 starter / extender / 手坑 / 终场怪兽
- ✅ 起手概率计算
- ❌ 不评价卡组好不好
- ❌ 不分析弱点或对局表现
- ❌ 不给建议

对局分析由 `deck-analyzer` skill 负责。

---

## 输入

读取工作目录中的 `deck_raw.json`，内容为 `get_cards_by_ydk` 和 `count-deck` 的原始返回数据。

---

## 解析流程

### Step 1 — 卡片分类

读取每张卡的效果文本，将每张卡归入以下类型。一张卡可以归入多个类型。

| 类型 | 定义 |
|---|---|
| `starter` | 能独立启动主引擎展开、或检索到关键卡的卡片。通常是 normal summon 后搜索、或一卡展开的魔法。 |
| `extender` | 无法独立启动展开，但能在 combo 进行中提供额外召唤点或资源的卡片。通常是从手牌/墓地特殊召唤的补点。 |
| `handtrap` | 在对手回合从手牌发动效果来干扰对手的卡片。 |
| `engine_core` | 属于某个 archetype 的核心卡，但不属于上面三类（例如 combo 中途的素材、终场怪兽的素材来源）。 |
| `endboard_monster` | 预期作为终场怪兽留场的额外卡组怪兽。 |
| `mid_material` | 主要用作中途 link/xyz/synchro 素材的额外卡组怪兽。 |
| `boardbreaker` | 主要用于突破对手盘面的卡片，例如冥王结界波、拮抗者、坏兽等。 |
| `floodgate` | 持续限制对手行动的永续魔陷或怪兽效果。 |
| `resource_denial` | 通过破坏对手手牌、卡组或墓地资源来压制对手的卡片，例如轰雷帝、暗黑界军神。注意：这类卡通常本身就是某些卡组的终场核心，而非通用工具。 |
| `utility` | 不属于以上类型的泛用功能卡。 |

### 关于 once_per_turn_type

读取每张卡的效果文本后，同时标注效果发动的频率限制：

| 值 | 定义 |
|---|---|
| `card_name_once` | 效果有卡名限制，同名卡在同一回合只能发动一次。多张同名的意义仅是提高起手率。 |
| `once_per_turn` | 效果有一回合一次限制，但没有卡名绑定。在特定条件下（例如不同状态、不同位置）同一张卡可以多次发动，或多张同名可各自发动。 |
| `unlimited` | 效果没有次数限制，理论上可以在同一回合内多次发动。塞勒涅、神圣魔皇后等靠计数器反复利用的卡属于此类。 |

`once_per_turn` 和 `unlimited` 的卡在资源堆叠上有特殊价值，deck-analyzer 会用这个信息分析资源上限。

### Step 2 — Engine 识别

将属于同一 archetype 或同一功能模块的卡片归组，形成 engine 列表。

对每个 engine 记录：
- 包含哪些卡片
- 在卡组中的角色（主引擎 / 副引擎 / 独立 package）
- 是否依赖 Normal Summon 启动
- 是否与其他 engine 争抢资源（例如同一额外怪物区位置、同一召唤权）

### Step 3 — 起手概率计算

统计 starter 总张数和手坑总张数，调用 `starter-probability`，计算：
- 起手至少一张 starter 的概率
- 起手至少一张手坑的概率

---

## 输出格式

输出合法 JSON，dump 到 `deck_parsed.json`。

```json
{
  // 卡组基本信息
  "deck_profile": {
    "archetype": "卡组主题名称，例如 '蛇眼' / '白噪音'",
    "deck_type": "卡组风格，枚举值：'combo' | 'control' | 'midrange' | 'turbo'",
    "difficulty": "操作难度，枚举值：'low' | 'medium' | 'high'"
  },

  // 卡组规模，来自 count-deck
  "deck_structure": {
    "main": 0,   // 主卡组张数
    "extra": 0,  // 额外卡组张数
    "side": 0    // 副卡组张数
  },

  // Engine 列表，每个 engine 一个对象
  "engines": [
    {
      "name": "engine 名称，例如 '蛇眼引擎'",
      "role": "primary | secondary | package  // primary=主引擎, secondary=副引擎, package=独立功能包",
      "cards": [
        {
          "name": "卡片名称",
          "copies": 0,        // 卡组中张数
          "types": [],        // 此卡的类型列表，从上方类型枚举中选取
          "once_per_turn_type": "card_name_once | once_per_turn | unlimited  // 效果频率限制，见上方说明",
          "effect_summary": "一句话说明效果或在 engine 中的作用"
        }
      ],
      "normal_summon_dependent": true,  // 是否依赖召唤权启动
      "resource_conflicts": "与哪个 engine 存在资源冲突，没有则填 null"
    }
  ],

  // Starter 汇总
  "starters": {
    "cards": [
      {
        "name": "卡片名称",
        "copies": 0,
        "once_per_turn_type": "card_name_once | once_per_turn | unlimited",
        "effect_summary": "一句话说明这张卡如何启动展开，例如 '召唤后搜索同名，触发连锁展开'"
      }
    ],
    "total_copies": 0  // 所有 starter 张数之和
  },

  // Extender 汇总
  "extenders": {
    "cards": [
      {
        "name": "卡片名称",
        "copies": 0,
        "once_per_turn_type": "card_name_once | once_per_turn | unlimited",
        "effect_summary": "一句话说明这张卡如何补充展开，例如 '墓地有同名时可从手牌特召'"
      }
    ],
    "total_copies": 0
  },

  // 手坑汇总
  "handtraps": {
    "cards": [
      {
        "name": "卡片名称",
        "copies": 0,
        "once_per_turn_type": "card_name_once | once_per_turn | unlimited",
        "effect_summary": "一句话说明效果，例如 '对手从卡组特召时，可以无效并破坏'"
      }
    ],
    "total_copies": 0
  },

  // Boardbreaker 和封锁卡汇总（非手坑类，用于突破对手盘面或限制对手行动）
  "boardbreakers_and_floodgates": {
    "cards": [
      {
        "name": "卡片名称",
        "copies": 0,
        "type": "boardbreaker | floodgate | resource_denial",
        "once_per_turn_type": "card_name_once | once_per_turn | unlimited",
        "effect_summary": "一句话说明效果"
      }
    ]
  },

  // 额外卡组怪兽列表
  "extradeck_monsters": [
    {
      "name": "卡片名称",
      "copies": 0,
      "role": "endboard_monster | mid_material | both  // 终场怪兽 / 中途素材 / 两者皆是",
      "effect_summary": "一句话说明在 combo 中的作用"
    }
  ],

  // 起手概率，来自 starter-probability
  "probability": {
    "starter_open_rate": 0.0,   // 0.0–1.0，起手至少一张 starter 的概率
    "handtrap_open_rate": 0.0   // 0.0–1.0，起手至少一张手坑的概率
  }
}
```
