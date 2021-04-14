import { Option, Select } from "@leafygreen-ui/select";
import { useEffect, useState } from "react";
import { Metric, SingleChart, Snapshot } from "SingleChart";
import { H3 } from "@leafygreen-ui/typography";
import "./Timechart.css";

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
        });
    };

    getData();
    const interval = setInterval(getData, 10000);
    return () => clearInterval(interval);
  }, [collection]);

  return (
    <div>
      <div className="top-row">
        <H3>{`Database: ${snapshots && snapshots[0] && snapshots[0].db}`}</H3>
        <Select
          label=""
          value={collection}
          onChange={(val: Metric) => setCollection(val)}
        >
          {collections()}
        </Select>
      </div>
      <div className="grid-container">
        <div className="top-left">
          <H3>Number of Documents</H3>
          <SingleChart
            Data={snapshots}
            Metric="num_docs"
            Label=""
            Conversion={noChange}
          ></SingleChart>
        </div>
        <div className="top-right">
          <H3>Data Size (uncompressed)</H3>
          <SingleChart
            Data={snapshots}
            Metric="size_data"
            Label="Size (kb)"
            Conversion={toKbyte}
          ></SingleChart>
        </div>
        <div className="bottom-left">
          <H3>Data Size on Disk</H3>
          <SingleChart
            Data={snapshots}
            Metric="size_storage"
            Label="Size (kb)"
            Conversion={toKbyte}
          ></SingleChart>
        </div>
        <div className="bottom-right">
          <H3>Index Size</H3>
          <SingleChart
            Data={snapshots}
            Metric="size_indices"
            Label="Size (kb)"
            Conversion={toKbyte}
          ></SingleChart>
        </div>
      </div>
    </div>
  );
};
