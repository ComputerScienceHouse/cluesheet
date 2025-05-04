import { cluesheetBackendEndpoint } from "@/lib/endpoint";
import styles from "./page.module.scss";
import { Container } from "reactstrap";
import FlipNumbers from "react-flip-numbers";
import PointCounter from "@/components/viewer/counter/PointCounter";
import { Clue } from "@/lib/types";
import ClueGroup from "@/components/viewer/cluegroup/ClueGroup";

export const metadata = {
  title: "View cluesheet",
  description: "View an mf cluesheet kerchoo",
};

export default async function CluesheetViewer({
  params,
}: {
  params: Promise<{ cluesheet_id: string }>;
}) {
  //const { cluesheet_id } = await params;
  //const cluesheet = await fetch(`${cluesheetBackendEndpoint}/${cluesheet_id}`);
  //console.log(`Got cluesheet object: ${await cluesheet.json()}`);

  const cluesheet = {
    title: "Chom skz",
    user_points: 69420,
    user_additional_score: ["-Your Bones", "+Adam Neulight's Car"],
    /*point_unit: "Kubernetes Clusters",*/
    clues: [
      {
        id: "my-sick-and-poggers-uuid",
        limit: 1,
        completions: 1,
        description: "Eat a whole can of beans",
        points: "+5",
        children: [
          {
            id: "my-sick-and-poggers-uuid-2",
            limit: 1,
            completions: 1,
            description: "With a fork",
            points: "+1",
            children: [],
          },
        ],
      },
    ],
  };

  return (
    <>
      <main>
        <Container>
          <div className={styles.header}>
            <h1>Cluesheet Name</h1>
            <div className={styles.pointCounter}>
              <PointCounter points={cluesheet.user_points}/>
              <h3>Points</h3>
            </div>
          </div>
          <div className={styles.clueList}>
            {cluesheet.clues.length > 0 &&
              cluesheet.clues.map((clue: Clue, index) => (
                <ClueGroup key={clue.id} clue={clue} />
              ))}
          </div>
        </Container>
      </main>
    </>
  );
}
