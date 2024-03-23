import { useState } from "react";
import ReviewsColumn from "./components/ReviewsColumn";
import AppSelect from "./components/AppSelect";
import { App } from "./types/Apps";

export interface Review {
  id: number;
  author: string;
  content: string;
  score: number;
  date: string;
}

const App: React.FC = () => {
  const [selectedApp, setSelectedApp] = useState<App | null>(null);

  const handleAppSelection = (app: App) => {
    setSelectedApp(app);
  };

  return (
    <div className="justify-top flex min-h-screen flex-col items-center gap-4">
      <h1 className="pt-12 text-5xl font-bold text-[#6a7edc]">Re:View</h1>
      <h2 className="text-xl font-light">
        What the people are saying about the top apps on iOS
      </h2>
      <AppSelect onAppSelect={handleAppSelection} />
      <ReviewsColumn selectedApp={selectedApp} />
    </div>
  );
};

export default App;
