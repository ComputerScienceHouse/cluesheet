import { Clue, UserClue } from "@/lib/types";

import styles from "./page.module.scss";
import { useState } from "react";
interface LineItemProps {
  clue: Clue;
}

export default function LineItem({ clue }: LineItemProps) {
  // Initialize state to track if the checkbox is checked
  const [isChecked, setIsChecked] = useState((clue as UserClue).completions !== undefined && Number((clue as UserClue).completions) > 0);

  // Handle checkbox change
  const handleClueChecked = () => {
    setIsChecked((prevChecked) => !prevChecked);
    // You can also perform any additional logic here, like updating the clue state
  };

  return (
    <>
      <div className={styles.lineItem}>
        { (clue as UserClue).completions !== undefined &&
        <input
          type="checkbox"
          className={styles.customCheckboxInput}
          onChange={handleClueChecked}
          checked={isChecked}
        />}
        <h5>
          {clue.points}: {clue.description}
        </h5>
      </div>
    </>
  );
}
