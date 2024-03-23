import { useState } from "react";
import ReviewsColumn from "./components/ReviewsColumn";
import AppSelect from "./components/AppSelect";
import { AppsListItem } from "./types/Apps";

export interface Review {
  id: number;
  author: string;
  content: string;
  score: number;
  date: string;
}

const App: React.FC = () => {
  const [selectedApp, setSelectedApp] = useState<AppsListItem | null>(null);

  const handleAppSelection = (appItem: AppsListItem) => {
    setSelectedApp(appItem);
  };

  return (
    <div className="justify-top flex min-h-screen flex-col items-center gap-4">
      <h1 className="pt-12 text-5xl font-bold text-[#6a7edc]">Re:View</h1>
      <h2 className="text-center text-xl font-light">
        What the people are saying about the top apps on iOS
      </h2>
      <AppSelect onAppSelect={handleAppSelection} />
      <ReviewsColumn selectedApp={selectedApp} />
    </div>
  );
};

export default App;
