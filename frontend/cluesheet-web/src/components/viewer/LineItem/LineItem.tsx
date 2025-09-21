import { Clue } from "@/lib/types";
import styles from "./page.module.scss";
import { useState } from "react";

interface LineItemProps {
  ancestors: Array<Clue>;
  clue: Clue;
  handleCheckboxChange: (clueId: string) => void;
}

export function LineItem({
  ancestors,
  clue,
  handleCheckboxChange,
}: LineItemProps) {
  const ancestorTags: string[] = ancestors.flatMap((ancestor) => ancestor.tags);

  const [isChecked, setIsChecked] = useState(false);

  const isStackable = clue.rule?.key === "stacks";

  return (
    <div className={styles.lineItem}>
      {isStackable && (
        <input
          type="number"
          value={clue.completions}
          onChange={(e) => handleCheckboxChange(clue.id, e.target.value)}
        />
      )}
      {!isStackable && (
        <input
          type="checkbox"
          value={clue.id}
          className={styles.customCheckboxInput}
          //checked={isChecked}
          onChange={() => handleCheckboxChange(clue.id)}
        />
      )}
      <div className={styles.points}>
        {clue.points} {clue.rule && `(${clue.rule.key})`}
      </div>
      <div className={styles.description}>{clue.description}</div>
      <div className={styles.tags}>
        {clue.tags
          .filter((tag) => !ancestorTags.includes(tag))
          .map((tag: string, index) => (
            <p className={styles.tag} key={index}>
              {tag}
            </p>
          ))}
      </div>
    </div>
  );
}
