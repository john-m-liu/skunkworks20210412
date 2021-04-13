import { Option, Select } from "@leafygreen-ui/select";
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

  temp_display_metric: number;
}

const toKbyte = (size) => size / 1024;
const noChange = (input) => input;

type Metric = keyof CollectionStats;

const metrics = {
  num_docs: {
    Label: "Number of Documents",
    Unit: "Number",
    Conversion: noChange,
  },
  num_indices: {
    Label: "Number of Indices",
    Unit: "Number",
    Conversion: noChange,
  },
  size_data: {
    Label: "Data Size (uncompressed)",
    Unit: "Size (kb)",
    Conversion: toKbyte,
  },
  size_storage: {
    Label: "Data Storage Size",
    Unit: "Size (kb)",
    Conversion: toKbyte,
  },
  size_indices: {
    Label: "Index Size",
    Unit: "Size (kb)",
    Conversion: toKbyte,
  },
  size_total: {
    Label: "Total Size",
    Unit: "Size (kb)",
    Conversion: toKbyte,
  },
};

interface Props {}

export const Timechart: React.FC<Props> = () => {
  let [snapshots, setSnapshots] = useState<Snapshot[]>([]);
  let [metric, setMetric] = useState<Metric>("size_data");

  const getData = () => {
    fetch("http://localhost:3001/snapshots")
      .then((res) => res.json())
      .then((snapshots) => {
        setSnapshots(snapshots);
        console.log(snapshots);
      });
  };

  const metricOptions = () => {
    let options = [];
    for (let metric in metrics) {
      console.log(metric);
      options.push(<Option value={metric}>{metrics[metric].Label}</Option>);
    }
    return options;
  };

  useEffect(() => {
    getData();
    const interval = setInterval(getData, 10000);
    return () => clearInterval(interval);
  }, []);

  snapshots = snapshots.map((s) => {
    s.stats.temp_display_metric = metrics[metric].Conversion(s.stats[metric]);
    return s;
  });

  return (
    <div>
      <Select
        label="Metric"
        value={metric}
        onChange={(val: Metric) => setMetric(val)}
      >
        {metricOptions()}
      </Select>
      <LineChart
        width={1500}
        height={800}
        data={snapshots}
        margin={{ top: 5, right: 20, bottom: 5, left: 0 }}
      >
        <Line
          type="monotone"
          dataKey="stats.temp_display_metric"
          isAnimationActive={false}
          stroke="#8884d8"
        />
        <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
        <XAxis
          dataKey="gathered_time_unix"
          tickFormatter={(unixTime) => moment.unix(unixTime).format("HH:mm:ss")}
          type="number"
          domain={["auto", "auto"]}
        />
        <YAxis
          dataKey="stats.temp_display_metric"
          type="number"
          domain={["auto", "auto"]}
          label={metrics[metric].Unit}
          width={200}
        />
        <Tooltip />
      </LineChart>
    </div>
  );
};
