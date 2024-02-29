import { useEffect, useRef, useState } from "react";
import "./App.css";
import { GetClipboardText, StartLLMStream } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime";
import { useHotkeys } from "react-hotkeys-hook";
import Markdown from "react-markdown";
import { Link, useNavigate } from "react-router-dom";
import { IconGear, IconPaste, IconTrashcan } from "./Icons";

function App() {
  const [previousConversation, setPreviousConversation] = useState(
    [] as string[]
  );
  const [clipboardText, setClipboardText] = useState("");
  const [promptText, setPromptText] = useState("");
  const [answer, setAnswer] = useState("");
  const [isWaiting, setWaiting] = useState(false);
  const scrollBottomMark = useRef<HTMLSpanElement>(null);

  useHotkeys(
    "mod+enter",
    () => {
      getLLMAnswer();
    },
    {
      enableOnFormTags: ["TEXTAREA"],
    }
  );

  useHotkeys(
    "mod+k",
    () => {
      if (clipboardText.length) {
        setClipboardText("");
      } else {
        getTextFromClipboard();
      }
    },
    {
      enableOnFormTags: ["TEXTAREA"],
    }
  );

  useHotkeys("escape", () => {
    if (answer.length) {
      setPreviousConversation((prev) => [...prev, promptText, answer]);
      setAnswer("");
      setPromptText("");
    }
  });

  const getTextFromClipboard = () => {
    GetClipboardText().then((text) => {
      setClipboardText(text);
    });
  };

  const getLLMAnswer = () => {
    setWaiting(true);
    let prompt = promptText;
    if (clipboardText.length) {
      prompt = `${promptText}\nContext:\n"""\n${clipboardText}\n"""`;
    }
    StartLLMStream(previousConversation, prompt);
  };

  useEffect(() => {
    EventsOn("answer-update", (content: string) => {
      setWaiting(false);
      setAnswer(content);
      if (scrollBottomMark.current) {
        scrollBottomMark.current.scrollIntoView();
      }
    });
  }, []);

  if (isWaiting || answer.length) {
    return (
      <div id="app">
        {isWaiting ? (
          <div className="px-2 py-2 rounded-md bg-white text-black focus:outline-none ring-2 ring-gray-300 focus:ring-blue-500 flex-grow overflow-y-auto flex justify-center items-center">
            Thinking...
          </div>
        ) : (
          <div className="px-2 py-2 rounded-md bg-white text-black focus:outline-none ring-2 ring-gray-300 focus:ring-blue-500 flex-grow overflow-y-auto markdown no-drag">
            <Markdown>{answer}</Markdown>
            <span id="scroll-bottom-mark" ref={scrollBottomMark}></span>
          </div>
        )}
        <button
          className="px-2 py-2 rounded-md bg-gray-200 hover:bg-gray-100 text-gray-800 hover:text-gray-700 flex gap-2 items-center"
          onClick={() => setAnswer("")}
        >
          <span className="bg-gray-300 px-2 py-1 text-xs rounded-md text-gray-500">
            Esc
          </span>
          <span className="flex-1">Go back</span>
          <span className="w-10 opacity-0"></span>
        </button>
      </div>
    );
  }

  return (
    <div id="app">
      <div className="flex items-center justify-between">
        <div>
          With context:{" "}
          <span className="text-xs rounded-md bg-gray-200 text-gray-800 px-2 py-1">
            {clipboardText.length
              ? `"${clipboardText.substring(0, 40)}..."`
              : "<empty>"}
          </span>{" "}
          {clipboardText.length ? (
            <button
              className="px-2 py-1 rounded-md bg-red-400 text-white text-xs"
              onClick={() => setClipboardText("")}
            >
              <div className="flex gap-1 items-center">
                <IconTrashcan />
                Clear
              </div>
            </button>
          ) : (
            <button
              className="px-2 py-1 rounded-md bg-green-500 text-white text-xs"
              onClick={() => getTextFromClipboard()}
            >
              <div className="flex gap-1 items-center">
                <IconPaste />
                Add clipboard
              </div>
            </button>
          )}
        </div>
        <Link to="/settings">
          <IconGear />
        </Link>
      </div>
      <textarea
        autoFocus={true}
        value={promptText}
        onChange={(e) => setPromptText(e.target.value)}
        placeholder="Enter a custom prompt here..."
        className="px-2 py-2 rounded-md bg-white text-black focus:outline-none ring-2 ring-gray-300 focus:ring-blue-500 resize-none flex-1 no-drag"
      ></textarea>
      <button
        className="px-2 py-2 rounded-md bg-blue-500 hover:bg-blue-400 text-white flex gap-2 items-center"
        onClick={() => getLLMAnswer()}
      >
        <span className="bg-blue-600 px-2 py-1 text-xs rounded-md text-blue-300">
          ⌘ ↵
        </span>
        <span className="flex-1">Submit</span>
        <span className="w-10 opacity-0"></span>
      </button>
    </div>
  );
}

export default App;
