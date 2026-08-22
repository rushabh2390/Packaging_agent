"use client";

import React, { useState, useRef, useEffect } from "react";
import axios from "axios";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import {
  Cpu,
  Send,
  Trash2,
  RotateCcw,
  Sliders,
  MessageSquare,
  Search,
  Bot,
  Box,
} from "lucide-react";

const API_BASE_URL = process.env.BACKEND_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export default function PackagingStandardHub() {
  const [prompt, setPrompt] = useState("");
  const [temperature, setTemperature] = useState(0.0);
  const [topK, setTopK] = useState(40);
  const [retrievalK, setRetrievalK] = useState(3);
  const [loading, setLoading] = useState(false);

  const [chatHistory, setChatHistory] = useState<
    Array<{ role: "user" | "assistant"; content: string }>
  >([]);

  const [inspectorContent, setInspectorContent] = useState<string>(
    "Any extracted tables, FEFCO style schemas, or packaging dimension layouts selected by the agent tools will stack right here."
  );

  // Reference for auto-scrolling
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  // Trigger smooth slide/scroll down on every new message or loading state change
  useEffect(() => {
    scrollToBottom();
  }, [chatHistory, loading]);

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!prompt.trim() || loading) return;

    const userMessage = prompt;
    setPrompt("");

    const updatedHistory = [
      ...chatHistory,
      { role: "user" as const, content: userMessage },
    ];
    setChatHistory(updatedHistory);
    setLoading(true);

    // Append empty assistant message to stream tokens into
    setChatHistory((prev) => [...prev, { role: "assistant", content: "" }]);

    try {
      const response = await fetch(`${API_BASE_URL}/pack`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          messages: updatedHistory,
          retrieval_k: retrievalK,
          temperature: temperature,
          top_k: topK,
        }),
      });

      if (!response.ok || !response.body) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let assistantText = "";

      while (true) {
        const { value, done } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        const lines = chunk.split("\n");

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const rawData = line.slice(6).trim();
            if (rawData === "[DONE]") break;

            try {
              const parsed = JSON.parse(rawData);

              // Update 3D Inspector View if spatial metadata arrives
              if (parsed.placements || parsed.bin_dimensions) {
                setInspectorContent(
                  `### Spatial Placement Metrics\n* **Fill Efficiency**: \`${parsed.fill_percentage || "85"}%\`\n* **Box Dimension**: \`${parsed.bin_dimensions?.join("×") || "320×220×150"}\` mm\n`
                );
              }

              // Append streaming AI text chunks
              if (parsed.ai_recommendation) {
                assistantText += parsed.ai_recommendation;
                setChatHistory((prev) => {
                  const newHistory = [...prev];
                  newHistory[newHistory.length - 1] = {
                    role: "assistant",
                    content: assistantText,
                  };
                  return newHistory;
                });
              }
            } catch (e) {
              // Ignore partial SSE chunk parse errors
            }
          }
        }
      }
    } catch (err: any) {
      const errorMessage = err.message || "Unable to connect to backend engine.";
      setChatHistory((prev) => [
        ...prev.slice(0, -1),
        {
          role: "assistant",
          content: `**Error:** ${errorMessage}. Please check if your Go service is running on \`${API_BASE_URL}\`.`,
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#0a0f18] text-slate-200 font-sans flex flex-col p-4 gap-4">
      {/* Header */}
      <header className="flex items-center justify-between py-2.5 px-4 bg-[#111827] border border-slate-800 rounded-lg">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-emerald-950 text-emerald-400 rounded-md border border-emerald-800/50">
            <Box className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-base font-semibold tracking-wide text-slate-100">
              3D Packaging Standard & Sizing Engine
            </h1>
            <p className="text-[11px] text-slate-400">
              Provide product details, material, and dimensions to get outer box dimensions and FEFCO styles.
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2 text-xs bg-slate-900 border border-slate-800 px-3 py-1 rounded-full text-slate-300">
          <Cpu className="w-3.5 h-3.5 text-emerald-400" />
          <span>Fiber Go Engine</span>
        </div>
      </header>

      {/* Main Workspace */}
      <main className="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-4">
        {/* Left Column: Active Conversation Space & Inputs */}
        <div className="lg:col-span-9 flex flex-col gap-4">

          {/* Active Conversation Space (Fixed Height with Auto-Slide/Scroll) */}
          <div className="h-[680px] bg-[#111827] border border-slate-800/80 rounded-xl p-4 flex flex-col">
            <div className="flex items-center gap-2 pb-3 mb-3 border-b border-slate-800 text-xs font-semibold uppercase tracking-wider text-slate-400 shrink-0">
              <MessageSquare className="w-4 h-4 text-emerald-400" />
              Active Conversation Space
            </div>

            {/* Scrollable Message History Area */}
            <div className="flex-1 overflow-y-auto pr-2 flex flex-col gap-4 scroll-smooth">
              {chatHistory.length === 0 ? (
                <div className="m-auto text-center text-xs text-slate-500">
                  <Bot className="w-8 h-8 text-slate-600 mx-auto mb-2" />
                  No active conversation. Provide your product specifications below to start sizing.
                </div>
              ) : (
                chatHistory.map((msg, index) => (
                  <div
                    key={index}
                    className={`p-4 rounded-xl text-xs leading-relaxed transition-all duration-300 animate-in fade-in slide-in-from-bottom-2 ${msg.role === "user"
                        ? "bg-slate-800 text-slate-200 self-end max-w-[80%]"
                        : "bg-[#0b0f19] border border-slate-800 text-slate-300 self-start w-full"
                      }`}
                  >
                    <p className="font-semibold text-[10px] uppercase mb-2 text-emerald-400 tracking-wider">
                      {msg.role === "user" ? "Product Request" : "Packaging Agent"}
                    </p>

                    <div className="prose prose-invert prose-xs max-w-none prose-p:leading-relaxed prose-pre:bg-slate-900 prose-table:border-slate-800 prose-th:border-slate-800 prose-td:border-slate-800">
                      <ReactMarkdown remarkPlugins={[remarkGfm]}>
                        {msg.content}
                      </ReactMarkdown>
                    </div>
                  </div>
                ))
              )}

              {loading && (
                <div className="p-3 bg-[#0b0f19] border border-slate-800 text-emerald-400 rounded-xl text-xs self-start flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-emerald-400 animate-ping" />
                  Calculating optimal box dimensions and FEFCO codes...
                </div>
              )}

              {/* Invisible anchor element to scroll down to */}
              <div ref={messagesEndRef} />
            </div>
          </div>

          {/* Input Box Always Positioned Below */}
          <div className="bg-[#111827] border border-slate-800/80 rounded-xl p-4">
            <div className="flex items-center justify-between pb-2">
              <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-slate-400">
                <Bot className="w-4 h-4 text-emerald-400" />
                Product & Packaging Input
              </div>
              <span className="text-[10px] text-slate-500">
                Format: Product Name | Material | Size (W×H×D mm)
              </span>
            </div>

            <form onSubmit={handleSendMessage} className="relative mt-2">
              <input
                type="text"
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="e.g., Glass bottles, fragile material, 300x200x120mm size — need FEFCO box code and outer dimensions"
                className="w-full bg-[#0b0f19] border border-slate-800 rounded-lg py-3 pl-4 pr-12 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:border-emerald-500/70 transition"
              />
              <button
                type="submit"
                disabled={loading}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 bg-emerald-600 hover:bg-emerald-500 disabled:bg-slate-700 text-white rounded-md transition"
              >
                <Send className="w-3.5 h-3.5" />
              </button>
            </form>
          </div>
        </div>

        {/* Right Sidebar: Agent Settings */}
        <aside className="lg:col-span-3 bg-[#111827] border border-slate-800/80 rounded-xl p-4 flex flex-col justify-between gap-6">
          <div className="flex flex-col gap-5">
            <div>
              <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-slate-300">
                <Sliders className="w-4 h-4 text-emerald-400" />
                Agent Settings
              </div>
              <p className="text-[10px] text-slate-500 mt-1">
                Engine: Go Packaging Service
              </p>
            </div>

            <div className="flex flex-col gap-4">
              <div className="flex flex-col gap-1.5">
                <div className="flex justify-between items-center text-xs">
                  <span className="text-slate-400">Retrieval K (Box Styles)</span>
                  <span className="text-emerald-400 font-mono font-bold">{retrievalK}</span>
                </div>
                <input
                  type="range"
                  min="1"
                  max="10"
                  value={retrievalK}
                  onChange={(e) => setRetrievalK(Number(e.target.value))}
                  className="w-full accent-emerald-500 bg-slate-800 h-1.5 rounded-lg cursor-pointer"
                />
              </div>

              <div className="flex flex-col gap-1.5">
                <div className="flex justify-between items-center text-xs">
                  <span className="text-slate-400">Generation Temperature</span>
                  <span className="text-emerald-400 font-mono font-bold">
                    {temperature.toFixed(2)}
                  </span>
                </div>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.05"
                  value={temperature}
                  onChange={(e) => setTemperature(Number(e.target.value))}
                  className="w-full accent-emerald-500 bg-slate-800 h-1.5 rounded-lg cursor-pointer"
                />
              </div>

              <div className="flex flex-col gap-1.5">
                <div className="flex justify-between items-center text-xs">
                  <span className="text-slate-400">Generation Top-K Window</span>
                  <span className="text-emerald-400 font-mono font-bold">{topK}</span>
                </div>
                <input
                  type="range"
                  min="1"
                  max="100"
                  value={topK}
                  onChange={(e) => setTopK(Number(e.target.value))}
                  className="w-full accent-emerald-500 bg-slate-800 h-1.5 rounded-lg cursor-pointer"
                />
              </div>
            </div>
          </div>

          <div className="flex flex-col gap-2 pt-4 border-t border-slate-800">
            <button
              onClick={() => setChatHistory([])}
              className="w-full py-2 px-3 bg-slate-800/80 hover:bg-slate-700/80 border border-slate-700 text-slate-300 rounded-lg text-xs font-medium flex items-center justify-center gap-2 transition"
            >
              <RotateCcw className="w-3.5 h-3.5" />
              Clear Conversation
            </button>
          </div>
        </aside>
      </main>
    </div>
  );
}