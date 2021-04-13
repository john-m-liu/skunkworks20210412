import moment from "moment";
import { useEffect, useState } from "react";
import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
} from "recharts";

interface Snapshot {
  db: string;
  collection: string;
  gathered_time: Date;
  gathered_time_unix: number;
  stats: CollectionStats;
}

interface CollectionStats {
  collection: string;
  num_docs: number;
  num_indices: number;
  size_data: number;
  size_storage: number;
  size_indices: number;
  size_total: number;

  size_total_kb: number;
  size_data_kb: number;
}

interface Props {}

export const Timechart: React.FC<Props> = () => {
  let [snapshots, setSnapshots] = useState<Snapshot[]>([]);
  useEffect(() => {
    fetch("http://localhost:3001/snapshots")
      .then((res) => res.json())
      .then((snapshots) => {
        setSnapshots(snapshots);
        console.log(snapshots);
      });
  }, []);

  snapshots = snapshots.map((s) => {
    s.stats.size_total_kb = toKbyte(s.stats.size_total);
    s.stats.size_data_kb = toKbyte(s.stats.size_data);
    return s;
  });

  return (
    <LineChart
      width={1500}
      height={800}
      data={snapshots}
      margin={{ top: 5, right: 20, bottom: 5, left: 0 }}
    >
      <Line type="monotone" dataKey="stats.size_data_kb" stroke="#8884d8" />
      <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
      <XAxis
        dataKey="gathered_time_unix"
        tickFormatter={(unixTime) => moment.unix(unixTime).format("HH:mm:ss")}
        type="number"
        domain={["auto", "auto"]}
      />
      <YAxis
        dataKey="stats.size_data_kb"
        type="number"
        domain={["auto", "auto"]}
        label="Size (kb)"
        width={200}
      />
      <Tooltip />
    </LineChart>
  );
};

const toKbyte = (size) => size / 1024;
