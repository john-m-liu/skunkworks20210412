import moment from "moment";
import React from "react";
import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
} from "recharts";

export type Snapshot = {
  db: string;
  collection: string;
  gathered_time: Date;
  gathered_time_unix: number;
  stats: CollectionStats;
};

export type CollectionStats = {
  collection: string;
  num_docs: number;
  num_indices: number;
  size_data: number;
  size_storage: number;
  size_indices: number;
  size_total: number;

  temp_display_metric: number;
};

export type Metric = keyof CollectionStats;

interface Props {
  Data: Snapshot[];
  Metric: Metric;
  Label: string;
  Conversion: (number) => number;
}

export const SingleChart: React.FC<Props> = (props) => {
  let data = [];
  props.Data.map((s) => {
    return data.push({
      x: s.gathered_time_unix,
      value: props.Conversion(s.stats[props.Metric]),
    });
  });

  return (
    <LineChart
      width={800}
      height={400}
      data={data}
      margin={{ top: 5, right: 0, bottom: 5, left: 0 }}
    >
      <Line
        type="monotone"
        dataKey="value"
        isAnimationActive={false}
        stroke="#8884d8"
      />
      <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
      <XAxis
        dataKey="x"
        tickFormatter={(unixTime) => moment.unix(unixTime).format("HH:mm:ss")}
        type="number"
        domain={["auto", "auto"]}
      />
      <YAxis
        dataKey="value"
        type="number"
        domain={["auto", "auto"]}
        label={props.Label}
        width={200}
      />
      <Tooltip
        isAnimationActive={false}
        labelFormatter={(value) => moment.unix(value).format("HH:mm:ss")}
      />
    </LineChart>
  );
};
