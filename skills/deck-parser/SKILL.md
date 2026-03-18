---
name: deck-parser
description: 游戏王卡组结构解析 skill。接收 CLI 原始卡片数据（deck_raw.csv），识别 engine、starter、extender、手坑，计算起手概率，输出结构化卡组数据 JSON。只做事实提取，不做对局判断。通常由 ygo-deck-analyze 调用，输入为 deck_raw.csv，输出为 deck_parsed.json。
---

# deck-parser

## 职责边界

这个 skill **只做事实提取**，不做判断。

- ✅ 卡片是什么、效果是什么、分几张
- ✅ 这张卡属于哪个 engine
- ✅ 这张卡是 starter / extender / 手坑 / 终场
- ✅ 起手概率计算
- ❌ 不评价卡组好不好
- ❌ 不分析弱点或对局表现
- ❌ 不给建议

对局分析由 `deck-analyzer` skill 负责。

---

## 输入

读取工作目录中的 `deck_raw.csv`，内容为 `get-cards-by-ydk` CLI 工具的原始返回数据（CSV格式）。

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
| `endboard` | 在终场留场、持续提供抗性/干扰/效果的卡片。**不限于额外卡组怪兽**，主卡组的效果怪兽、陷阱怪兽、永续魔法/陷阱均适用。判断标准：这张卡在终场是否还在场上发挥持续作用？ |
| `mid_material` | 主要用作中途 link/xyz/synchro 素材的额外卡组怪兽。 |
| `boardbreaker` | 主要用于突破对手盘面的卡片，例如冥王结界波、拮抗者、坏兽等。 |
| `floodgate` | 持续限制对手行动的永续魔陷或怪兽效果。 |
| `resource_denial` | 通过破坏对手手牌、卡组或墓地资源来压制对手的卡片，例如轰雷帝、暗黑界军神。注意：这类卡通常本身就是某些卡组的终场核心，而非通用工具。 |
| `utility` | 不属于以上类型的泛用功能卡。 |

#### 关于 `endboard` 的多标签说明

`extender` 和 `endboard` 对同一张卡经常同时成立，两个标签不互斥。
典型情况：一张卡在展开中作为补点（extender），但留场后持续提供抗性或干扰效果（endboard）。
**不要因为一张卡是 extender 就排除它成为 endboard 的可能性，必须独立判断它的终场价值。**

#### 关于陷阱怪兽（Trap Monsters）

陷阱怪兽具有双重身份，分类时注意：
- 在手牌/盖放时是**陷阱**，在场上特召后是**怪兽**。

---

### 关于 once_per_turn_type

读取每张卡的效果文本后，同时标注效果发动的频率限制：

| 值 | 定义 |
|---|---|
| `card_name_once` | 效果有卡名限制，同名卡在同一回合只能发动一次。多张同名的意义仅是提高起手率。 |
| `once_per_turn` | 效果有一回合一次限制，但没有卡名绑定。在特定条件下（例如不同状态、不同位置）同一张卡可以多次发动，或多张同名可各自发动。 |
| `unlimited` | 效果没有次数限制，理论上可以在同一回合内多次发动。塞勒涅、神圣魔皇后等靠计数器反复利用的卡属于此类。 |

`once_per_turn` 和 `unlimited` 的卡在资源堆叠上有特殊价值，deck-analyzer 会用这个信息分析资源上限。

---

### Step 2 — Engine 识别

将属于同一 archetype 或同一功能模块的卡片归组，形成 engine 列表。

Engine 归组按以下优先级判断，**高优先级规则命中时，即使 archetype 字段指向其他分组，也以高优先级规则为准**：

**⑤ 召唤/发动条件依赖**（最高优先级）
如果卡片 A 的召唤条件或效果发动条件中，要求场上/墓地存在某张具名卡 B 或某个 archetype 的卡片，则 A 与 B / 该 archetype 存在强依赖，归入同一 engine。
典型特征：
- "自己场上有「X」存在的场合才能特殊召唤"
- "自己场上有「X」archetype 的卡存在时才能发动"

即使 A 有独立的 archetype 字段，也应优先按依赖关系归组。

**② 效果文本卡名直接引用**（高优先级）
如果卡片 A 的效果文本中出现了卡片 B 的完整卡名（通常以「」括起），则 A 与 B 关联，归入同一 engine。即使 A 没有任何 archetype 字段也适用。

**③ 效果文本 Archetype 名引用**（中优先级）
如果卡片 A 的效果文本中出现了某个 archetype 名称，且卡组中存在该 archetype 的卡片 B，则 A 与 B 关联，归入同一 engine。

**④ 功能性产出依赖**（中优先级）
如果一张场地/永续魔法/永续陷阱的效果能直接产出（搜索/特召/盖放）某个 archetype 的卡片，则该卡与该 archetype 归为同一 engine，不论自身 archetype 字段是否为空。

**① Archetype 字段匹配**（兜底）
archetype 字段相同的卡归入同一 engine。这是最弱的归组依据，当上方任一规则命中时，以上方规则为准。

对每个 engine 记录：
- 包含哪些卡片（含通过文本关联加入的）
- 在卡组中的角色（主引擎 / 副引擎 / 独立 package）
- 是否依赖 Normal Summon 启动
- 是否与其他 engine 争抢资源（例如同一额外怪物区位置、同一召唤权）

---

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
          "copies": 0,
          "types": [],        // 此卡的类型列表，从上方类型枚举中选取，endboard 和 extender 可同时存在
          "once_per_turn_type": "card_name_once | once_per_turn | unlimited",
          "join_reason": "archetype字段匹配 | 效果引用「X」| archetype名引用「X」| 功能性产出「X」| 召唤/发动条件依赖「X」  // 说明这张卡为何归入本 engine，archetype字段匹配时可省略",
          "endboard_role": "这张卡在终场的持续贡献，例如'为全场陷阱怪兽提供对象抗性；被破坏时触发联动'。仅在 types 含 endboard 时填写，否则填 null",
          "effect_summary": "一句话说明效果或在 engine 中的作用"
        }
      ],
      "normal_summon_dependent": true,
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
        "effect_summary": "一句话说明这张卡如何启动展开"
      }
    ],
    "total_copies": 0
  },

  // Extender 汇总
  "extenders": {
    "cards": [
      {
        "name": "卡片名称",
        "copies": 0,
        "once_per_turn_type": "card_name_once | once_per_turn | unlimited",
        "effect_summary": "一句话说明这张卡如何补充展开"
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
        "effect_summary": "一句话说明效果"
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
      "role": "endboard | mid_material | both  // 终场 / 中途素材 / 两者皆是",
      "effect_summary": "一句话说明在 combo 中的作用"
    }
  ],

  // 主卡组终场核心列表
  // 主卡组中 types 含 endboard 的卡片，在此单独汇总，便于 deck-analyzer 直接定位
  "maindeck_endboard": [
    {
      "name": "卡片名称",
      "copies": 0,
      "card_type": "effect_monster | trap_monster | continuous_spell | continuous_trap | field_spell | other",
      "once_per_turn_type": "card_name_once | once_per_turn | unlimited",
      "endboard_role": "这张卡在终场的持续贡献",
      "effect_summary": "一句话说明效果"
    }
  ],

  // 起手概率，来自 starter-probability
  "probability": {
    "starter_open_rate": 0.0,   // 0.0–1.0，起手至少一张 starter 的概率
    "handtrap_open_rate": 0.0   // 0.0–1.0，起手至少一张手坑的概率
  }
}
```
