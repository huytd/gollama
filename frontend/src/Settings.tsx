import { useEffect, useState } from "react";
import { GetLLMConfig, SetLLMConfig } from "../wailsjs/go/main/App";
import { useNavigate } from "react-router-dom";

function Settings() {
  const navigation = useNavigate();
  const [config, setConfig] = useState<{
    APIUrl?: string;
    APIModel?: string;
    APIKey?: string;
  }>({});

  useEffect(() => {
    GetLLMConfig().then((configData) => {
      try {
        const parsedConfig = JSON.parse(configData);
        setConfig(parsedConfig);
      } catch (e) {
        console.error(e);
      }
    });
  }, []);

  const saveConfig = () => {
    const apiUrl = config.APIUrl ?? "http://localhost:11343/v1";
    const apiKey = config.APIKey ?? "ollama";
    const apiModel = config.APIModel ?? "mistral:latest";
    SetLLMConfig(apiUrl, apiKey, apiModel);
    navigation("/");
  };

  return (
    <div className="text-black w-full h-full">
      <h1 className="font-bold text-lg border-b px-4 py-3 bg-gray-100">
        Settings
      </h1>
      <div className="flex h-full">
        <ul className="w-4/12 border-r p-4">
          <li className="py-1 px-4 rounded-md bg-blue-500 text-white mb-1 cursor-default">
            LLM Config
          </li>
          <li className="py-1 px-4 rounded-md hover:bg-blue-500 hover:text-white mb-1 cursor-default">
            Prompt Library
          </li>
        </ul>
        <div className="w-8/12 p-4">
          <div className="mb-3 w-full">
            <label htmlFor="apiUrl" className="block font-medium text-gray-700">
              API URL
            </label>
            <input
              type="text"
              id="apiUrl"
              name="apiUrl"
              className="mt-1 px-2 py-1 border border-gray-300 rounded-md w-full"
              defaultValue={config.APIUrl ?? ""}
              onChange={(e) => {
                const value = e.target.value;
                setConfig((prev) => ({ ...prev, APIUrl: value }));
              }}
            />
          </div>
          <div className="mb-3 w-full">
            <label htmlFor="apiKey" className="block font-medium text-gray-700">
              API Key
            </label>
            <input
              type="text"
              id="apiKey"
              name="apiKey"
              className="mt-1 px-2 py-1 border border-gray-300 rounded-md w-full"
              defaultValue={config.APIKey ?? ""}
              onChange={(e) => {
                const value = e.target.value;
                setConfig((prev) => ({ ...prev, APIKey: value }));
              }}
            />
          </div>
          <div className="mb-3 w-full">
            <label
              htmlFor="selectedModel"
              className="block font-medium text-gray-700"
            >
              Default Model
            </label>
            <input
              type="text"
              id="apiModel"
              name="apiModel"
              className="mt-1 px-2 py-1 border border-gray-300 rounded-md w-full"
              defaultValue={config.APIModel ?? ""}
              onChange={(e) => {
                const value = e.target.value;
                setConfig((prev) => ({ ...prev, APIModel: value }));
              }}
            />
          </div>
          <button
            className="px-4 py-1 text-white bg-blue-500 border-2 border-blue-500 rounded-md"
            onClick={saveConfig}
          >
            Save
          </button>
        </div>
      </div>
    </div>
  );
}

export default Settings;
