"use client";
import FlipNumbers from "react-flip-numbers";
import styles from "./page.module.scss";

interface PointCounterProps {
  points: number;
}

export default function PointCounter({ points }: PointCounterProps) {
  return (
    <div className={styles.flipNumbers}>
      <FlipNumbers
        height={40}
        width={40}
        color="black"
        background="white"
        play
        perspective={400}
        numbers={String(points)}
        numberClassName={styles.number}
      />
    </div>
  );
};
