import React, { useEffect, useState } from "react";
import { App } from "../types/Apps";
import { getTopApps } from "../api/getTopApps";

interface AppSelectProps {
  onAppSelect: (app: App) => void;
}

const AppSelect: React.FC<AppSelectProps> = ({ onAppSelect }) => {
  const [selection, setSelection] = useState<App | null>(null);
  const [apps, setApps] = useState<App[]>([]);

  useEffect(() => {
    const fetchApps = async () => {
      try {
        const response = await getTopApps();
        const data: App[] = [];
        // hardcode the test app
        data.push({
          Id: "595068606",
          Name: "Test App",
          IconUrl: "",
        });
        // merge the data with the response
        data.push(...response.Apps);
        // set the list of apps and select the first one
        setApps(data);
        onAppSelect(data[0]);
      } catch (error) {
        console.error("Failed to load apps: ", error);
      }
    };

    fetchApps();
  }, []);

  const handleSelectionChange = (
    event: React.ChangeEvent<HTMLSelectElement>,
  ) => {
    const selectedApp = apps.find((app) => app.Id === event.target.value);
    onAppSelect(selectedApp || apps[0]);
    setSelection(selectedApp || apps[0]);
  };

  return (
    <div className="border-g flex w-1/2 flex-row gap-2 p-4">
      <label htmlFor="app-select"></label>
      <select
        id="app-select"
        value={selection?.Id || ""}
        onChange={handleSelectionChange}
        className="w-full rounded-lg p-2 shadow-sm"
      >
        {apps.length === 0 && <option value=""></option>}
        {apps.map((app) => (
          <option key={app.Id} value={app.Id}>
            {app.Name}
          </option>
        ))}
      </select>
    </div>
  );
};

export default AppSelect;
