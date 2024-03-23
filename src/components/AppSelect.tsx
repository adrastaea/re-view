import React, { useEffect, useState } from "react";
import { AppsListItem } from "../types/Apps";
import { getTopApps } from "../api/getTopApps";

interface AppSelectProps {
  onAppSelect: (appItem: AppsListItem) => void;
}

const AppSelect: React.FC<AppSelectProps> = ({ onAppSelect }) => {
  const [selection, setSelection] = useState<AppsListItem | null>(null);
  const [appsList, setAppsList] = useState<AppsListItem[]>([]);

  useEffect(() => {
    const fetchApps = async () => {
      try {
        const response = await getTopApps();
        const data: AppsListItem[] = [];
        // hardcode the test app
        data.push({
          Id: "595068606",
          Name: "Test App",
          IconUrl: "",
        } as AppsListItem);
        // merge the data with the response
        data.push(...response.Apps);
        // set the list of apps and select the first one
        setAppsList(data);
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
    const selectedApp = appsList.find(
      (appItem) => appItem.Id === event.target.value,
    );
    onAppSelect(selectedApp || appsList[0]);
    setSelection(selectedApp || appsList[0]);
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
        {appsList.length === 0 && <option value=""></option>}
        {appsList.map((appItem) => (
          <option key={appItem.Id} value={appItem.Id}>
            {appItem.Name}
          </option>
        ))}
      </select>
    </div>
  );
};

export default AppSelect;
