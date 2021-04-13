import { Option, Select } from "@leafygreen-ui/select";
import { useEffect, useState } from "react";
import { Metric, SingleChart, Snapshot } from "SingleChart";

const toKbyte = (size) => size / 1024;
const noChange = (input) => input;

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

  return (
    <div>
      <Select
        label="Metric"
        value={metric}
        onChange={(val: Metric) => setMetric(val)}
      >
        {metricOptions()}
      </Select>
      <SingleChart
        Data={snapshots}
        Metric="num_docs"
        Label="Number"
        Conversion={noChange}
      ></SingleChart>
      <SingleChart
        Data={snapshots}
        Metric="size_data"
        Label="Size (kb)"
        Conversion={toKbyte}
      ></SingleChart>
      <SingleChart
        Data={snapshots}
        Metric="size_storage"
        Label="Size (kb)"
        Conversion={toKbyte}
      ></SingleChart>
      <SingleChart
        Data={snapshots}
        Metric="size_indices"
        Label="Size (kb)"
        Conversion={toKbyte}
      ></SingleChart>
    </div>
  );
};
