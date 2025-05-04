import { cluesheetBackendEndpoint } from "@/lib/endpoint";
import styles from "./page.module.scss";
import { Container } from "reactstrap";
import FlipNumbers from "react-flip-numbers";
import PointCounter from "@/components/viewer/counter/PointCounter";
import { Clue } from "@/lib/types";
import ClueGroup from "@/components/viewer/cluegroup/ClueGroup";
import { mockCluesheet } from "../../../../../../tests/lib/data";

export const metadata = {
  title: "View cluesheet",
  description: "View an mf cluesheet kerchoo",
};

export default async function CluesheetEditor({
  params,
}: {
  params: Promise<{ cluesheet_id: string; keycloak_uid: string }>;
}) {
  //const { cluesheet_id } = await params;
  //const cluesheet = await fetch(`${cluesheetBackendEndpoint}/api/v1/cluesheet/${cluesheet_id}`);
  //console.log(`Got cluesheet object: ${await cluesheet.json()}`);

  return (
    <>
      <main>
        <Container>
          <div className={styles.header}>
            <h1>{mockCluesheet.title}</h1>
            <div className={styles.pointCounter}>
              <PointCounter points={mockCluesheet.user_points} />
              <h3>Points</h3>
            </div>
          </div>
          <div className={styles.clueList}>
            {mockCluesheet.clues.length > 0 &&
              mockCluesheet.clues.map((clue: Clue) => (
                <ClueGroup key={`${clue.id}`} clue={clue} />
              ))}
          </div>
        </Container>
      </main>
    </>
  );
}
