# Yu-Gi-Oh! AI Tools 测试工作目录

这个目录演示如何使用我准备的这些 AI 工具。

## 安装步骤

1. **安装 ygo-db-cli**：
   ```bash
   go install github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/cmd/ygo-db-cli
   ```

2. **安装支持 skills 的 AI Agent**：
   本示例使用 Gemini CLI，因为它支持 skills 功能且 context window 达到 1M。

3. **安装 skills**：
   ```bash
   gemini skills install https://github.com/elmhuangyu/yu-gi-oh-ai-tools.git --path skills/ --scope workspace
   ```

完成以上步骤后，你就可以开始使用这些 AI 工具来分析你的卡组了！
