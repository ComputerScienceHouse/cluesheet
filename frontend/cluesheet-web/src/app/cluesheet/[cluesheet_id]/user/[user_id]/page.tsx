import { cluesheetBackendEndpoint } from "@/lib/endpoint";
import styles from "./page.module.scss";
import { Container } from "reactstrap";
import FlipNumbers from "react-flip-numbers";
import PointCounter from "@/components/viewer/counter/PointCounter";
import { Clue } from "@/lib/types";
import { mockUserCluesheet } from "../../../../../../tests/lib/data";
import { ClueList } from "@/components/viewer/ClueList/ClueList";
import CluesheetForm from "@/components/form/CluesheetForm/CluesheetForm";

export const metadata = {
  title: "View cluesheet",
  description: "Edit an mf cluesheet kerchoo",
};

export default async function CluesheetEditor({
  params,
}: {
  params: Promise<{ cluesheet_id: string; keycloak_uid: string }>;
}) {
  const { cluesheet_id, keycloak_uid } = await params;
  //const cluesheet = await fetch(`${cluesheetBackendEndpoint}/api/v1/cluesheet/${cluesheet_id}/user/${keycloak_uid}`);
  //console.log(`Got cluesheet object: ${await cluesheet.json()}`);

  const cluesheet = mockUserCluesheet;

  return (
    <>
      <main>
        <Container>
          <CluesheetForm cluesheet_id={cluesheet_id} keycloak_uid={keycloak_uid}/>
        </Container>
      </main>
    </>
  );
}
