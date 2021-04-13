import { Option, Select } from "@leafygreen-ui/select";
import { useEffect, useState } from "react";
import { Metric, SingleChart, Snapshot } from "SingleChart";

const toKbyte = (size) => size / 1024;
const noChange = (input) => input;

// const metrics = {
//   num_docs: {
//     Label: "Number of Documents",
//     Unit: "Number",
//     Conversion: noChange,
//   },
//   num_indices: {
//     Label: "Number of Indices",
//     Unit: "Number",
//     Conversion: noChange,
//   },
//   size_data: {
//     Label: "Data Size (uncompressed)",
//     Unit: "Size (kb)",
//     Conversion: toKbyte,
//   },
//   size_storage: {
//     Label: "Data Storage Size",
//     Unit: "Size (kb)",
//     Conversion: toKbyte,
//   },
//   size_indices: {
//     Label: "Index Size",
//     Unit: "Size (kb)",
//     Conversion: toKbyte,
//   },
//   size_total: {
//     Label: "Total Size",
//     Unit: "Size (kb)",
//     Conversion: toKbyte,
//   },
// };

interface Props {}

export const Timechart: React.FC<Props> = () => {
  let [snapshots, setSnapshots] = useState<Snapshot[]>([]);
  let [collection, setCollection] = useState<string>("collection0");

  const collections = () => {
    let c = [];
    for (let i = 0; i < 10; i++) {
      c.push(<Option value={`collection${i}`}>{`collection${i}`}</Option>);
    }
    return c;
  };

  useEffect(() => {
    const getData = () => {
      fetch(`http://localhost:3001/snapshots?collection=${collection}`)
        .then((res) => res.json())
        .then((snapshots) => {
          setSnapshots(snapshots);
          console.log(snapshots);
        });
    };

    getData();
    const interval = setInterval(getData, 10000);
    return () => clearInterval(interval);
  }, [collection]);

  return (
    <div>
      <Select
        label="Collection"
        value={collection}
        onChange={(val: Metric) => setCollection(val)}
      >
        {collections()}
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
