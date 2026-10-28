# PGR301 Eksamen 2026 — wishif.ai

<img alt="Wishifai" src="media/wishifai_banner.png">

Wishifai — en Oslo-basert startup — skal en gang for alle løse problemet med julegave- og bursdagsønsker, og strevet med å finne riktige gaver til familie, barn og venner.
Bare last opp et bilde av en stue, et barnerom eller garasjen til bestefar — Wishifai finner noe som passer. Prosjektet er naturligvis "vibe-kodet" av
solo-grunderen, og han trenger nå sårt å få orden på utviklingsprosessen før han lanserer i desember 2026.

Han har selv skrevet en Spring Boot-applikasjon som tar imot et bilde, laster det opp til S3, og bruker AWS Rekognition til å finne ut
hva slags objekter som finnes på bildet. Dette mates inn i et AWS Bedrock-prompt, som igjen genererer en liste med mulige gaveønsker som kan passe.

Din oppgave er å gjøre løsningen produksjonsklar: containerisere, deploye med Terraform, overvåke og sikre.

### Litt om Rekognition og Bedrock

**AWS Rekognition** er en ferdigtrent "computer vision"-tjeneste. Du sender inn et bilde, og får tilbake en liste med labels — f.eks. `Sofa (98%)`, `Lamp (87%)`, `Child (72%)` — pluss tekst, ansikter og objekter den gjenkjenner. Wishifai bruker labelene som input til neste steg.

**AWS Bedrock** er AWS sin administrerte tjeneste for store språkmodeller (Claude, Llama, Titan m.fl.). Wishifai bygger en prompt av Rekognition-labelene ("foreslå gaver basert på disse objektene i rommet …") og lar en LLM via Bedrock generere selve gavelisten. Bedrock er altså det som gjør tjenesten "AI-drevet" — og samtidig ikke-deterministisk, kostnadsvariabel og vanskelig å teste tradisjonelt. Det er temaer dere kommer tilbake til i oppgave 3, 4 og 5.


## Levering

- Opprett et **nytt** GitHub-repository (ikke fork). Kopier inn innholdet herfra.
- Besvar oppgavene i en fil du kaller `README_SVAR.md` i ditt eget repo.
- Lever repo-URL som PDF/TXT i Wiseflow.

## Utviklingsmiljø

Hele eksamen er løsbar i et **GitHub Codespace**, slik dere er vant til fra lab-ene.

Dette repoet inneholder en ferdig `.devcontainer/`-konfigurasjon med alle verktøyene som trengs forhåndsinstallert: Java 21, Maven, Docker, AWS CLI og Terraform. Hvis du vil bruke Codespaces kan du kopiere filen fra .devcontainer kappen i dette repoet til ditt eget.

Dere må selv konfigurere Codespaces med AWS nøkler som vanlig.

## Oppgaver

Totalt 100 poeng. Flere oppgaver gir rom for egne valg — begrunn valgene dine i `README_SVAR.md`.

Før du går videre bør du teste at Spring Boot-applikasjonen fungerer som forventet.

Start applikasjonen, og send inn et bilde med `curl` fra en annen terminal (fra repo-rot). Eksempelet under fungerer på Mac, Linux og i Codespaces:

```
curl -F "image=@samples/sample_image.png" http://localhost:8080/api/analyze
```
Som respons på et bilde får du et JSON objekt som inneholder gaveønsker og begrunnelse for valgene som er gjort

### Oppgave 1 — Container og CI (20p)

I denne oppgaven skal dere gjøre kodebasen klar for et større utviklingsteam, og sørge for at gode DevOps-praksiser blir fulgt.


#### 1A — Containerisér Spring Boot-applikasjonen

Dere skal lage en multi-stage `Dockerfile` for Spring Boot-applikasjonen. Grunderen har hittil deployet ved å zippe JAR-filen og scp-e den opp til en EC2-instans fra MacBooken sin. Det er dagens-gang-i-staur-territorium — vi skal gjøre dette ordentlig.

Multi-stage build betyr at build-verktøyene (Maven, JDK, kildekode) ikke henger med inn i det endelige imaget. Resultatet er mindre, raskere, og eksponerer færre avhengigheter enn nødvendig.

Leveranser: `Dockerfile`, `docker run …`-kommando sensor kan kjøre for å starte en container som lytter på localhost:9999

#### 1B — CI-pipeline

Dere skal sette opp en GitHub Actions-workflow som bygger Wishifai Spring Boot-containeren og publiserer den til **ECR**. ECR brukes senere i oppgave 2 når tjenesten skal kjøre på AWS; Docker Hub gjør det enkelt for sensor å trekke og teste imaget lokalt.

Workflowen skal **kun kjøre ved push til `main`**.

Dere velger selv en fornuftig tagging-strategi. `latest` alene er ikke tilstrekkelig — commit-SHA, versjon, eller tilsvarende må være med slik at det er mulig å spore hva som faktisk kjører hvor. Begrunn valget kort i `README_SVAR.md`.

ECR-repoet kan opprettes manuelt. Secrets skal ligge som GitHub Secrets, ikke i workflow-filen.

Leveranser: workflow-fil, lenke til grønn kjøring, begrunnelse av tagging-strategi i `README_SVAR.md`.

### Oppgave 2 — ECS Express med Terraform (25p)

Grunderen har så langt kjørt applikasjonen på en EC2-instans han manuelt setter opp i konsollet. Når én utvikler blir fem, holder ikke dét lenger — infrastrukturen må være kode, revisjonskontrollert, og mulig å gjenskape fra scratch.

I denne oppgaven skal dere flytte tjenesten over på **ECS Express**, Amazons forenklede container-tjeneste som gir dere "serverløse" containere og lastbalansering uten alt det manuelle oppsettet av nettverk og infrastrultur.

Nyttige referanser:
- AWS ECS Express Mode (konsepter, defaults, hva som opprettes): <https://docs.aws.amazon.com/AmazonECS/latest/developerguide/express-service-overview.html>
- Terraform-ressursen: <https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_express_gateway_service>

#### 2A — Infrastruktur for tjenesten

All Terraform-kode skal ligge i en mappe `infra/` i repoet. Tjenesten skal eksponere et offentlig endepunkt sensor kan curle. Bruk variabler for det som skal kunne endres (region, kandidatnummer, image-tag), og eksponér nyttige verdier via `outputs` — minimum selve URL-en.

**IAM-roller:** For å unngå rolleksplosjon i klassens AWS-konto, bruk de eksisterende rollene nedenfor i stedet for å opprette nye. Begge finnes allerede i kontoen og kan refereres direkte:

- **Execution role:** `arn:aws:iam::244530008913:role/ecsTaskExecutionRole`
- **Infrastructure role:** `arn:aws:iam::244530008913:role/service-role/ecsInfrastructureRoleForExpressServices`

For at tjenesten faktisk skal kunne kalle Rekognition, Bedrock og S3 trenger task-en en egen **task role** (ikke task execution role). Opprett én slik rolle i Terraform, suffikset med kandidatnummer, med nødvendige tillatelser.

Leveranser: Terraform-kode i `infra/`, offentlig URL til tjenesten.

#### 2B — Remote state

Terraform-state skal lagres i en **S3 backend**, ikke lokalt. Dette er grunnlaget for at flere utviklere kan jobbe mot samme infrastruktur uten å overskrive hverandre.

Bruk den delte bucketen **`pgr301-terraform-state`** i klassens AWS-konto. Velg en nøkkel som gir unikhet pr. kandidat, f.eks. `wishifai/<kandidatnummer>/terraform.tfstate`. Dokumentér valgt nøkkel i `README_SVAR.md`.

Leveranser: `backend`-konfigurasjon i Terraform-koden som peker på `pgr301-terraform-state`.

### Oppgave 3 — Observability (20p)

Grunderen oppdaget nylig at tjenesten hadde vært nede i tre dager før en kunde ringte for å spørre. Ingen varslet ham, og det eneste han har å gå på i ettertid er `System.out.println`-logger. Vi må kunne se hvordan systemet oppfører seg uten å ringe det.

#### 3A — Metrikker i applikasjonen

Dere skal instrumentere Spring Boot-applikasjonen med **Micrometer** og eksportere metrikkene til CloudWatch.

Applikasjonen har allerede nødvendige biblioteker på classpath (`micrometer-registry-cloudwatch2` + AWS CloudWatch SDK) og en ferdig `MetricsConfig`-klasse som wirer opp `CloudWatchMeterRegistry`. Dere trenger altså ikke sette opp selve CloudWatch-koblingen — fokus skal være på å *bruke* Micrometer til å måle det som er relevant.

Som et minimum skal følgende måles og publiseres som egne metrikker:

- Tiden det tar å kjøre **Rekognition**-kallet
- Tiden det tar å kjøre **Bedrock**-kallet
- Minst én **business-metrikk** — noe som sier noe om bruken av tjenesten, ikke bare teknisk helse. Eksempler: antall gaveforslag generert, antall bilder analysert pr. kategori (stue/barnerom/garasje), gjennomsnittlig antall forslag pr. prompt. Dere velger selv — begrunn valget.

Hvorfor Rekognition og Bedrock separat: når brukerne klager på at tjenesten er treg, må dere kunne svare på *hvor* tiden går. Et felles "latency total" er lite verdt når det er to eksterne tjenester som gjør grovjobben.

I tillegg står dere fritt til å legge til flere metrikker (`Counter`, `Gauge`, `DistributionSummary`) der det gir reell innsikt.

Beskriv kort i `README_SVAR.md` hvilke metrikker dere valgte og hvorfor.

#### 3B — Dashboard, alarm og Terraform-modul

Dashboard og alarm skal defineres som kode i `infra/`. Dashbordet skal gi et helhetlig bilde av tjenesten og skal vise en **miks av tre typer metrikker**:

- **Business:** minst én av dine egne business-metrikker fra 3A
- **Infrastruktur:** metrikker fra ECS Express (f.eks. CPU-utilization, memory-utilization, running task count)
- **Applikasjon:** metrikker eksponert via Spring Boot Actuator (f.eks. JVM memory, HTTP request-latency, GC-tid)

Minst én CloudWatch Alarm skal være knyttet til en av dine egne metrikker.

I denne oppgaven skal det også ligge minst **én egen Terraform-modul** (f.eks. `modules/dashboard/`) som kalles fra `main.tf`. Hensikten er å øve på gjenbrukbar infrastruktur, ikke å ha moduler for modulens egen skyld.

Leveranser: Terraform-kode, screenshot av dashboard med reelle data, kort begrunnelse av metrikkene.

### Oppgave 4 — Responsible AI (20p)

Grunderen fikk nylig en e-post fra en bruker som hadde lastet opp bilder av barnebarna sine. Tjenesten foreslo en gaveliste, men svaret inneholdt også detaljerte observasjoner om barna — estimert alder, antatt kjønn, utseende — rett fra det Rekognition hadde funnet. Grunderen svarte "beklager, vi tar det til etterretning" og gjør absolutt ingenting.

Vi skal gjøre noe.

#### 4A — Guardrail som kode

Dere skal konfigurere en **AWS Bedrock Guardrail** med Terraform. Policyen skal minst dekke PII (navn, adresser, kontonummer, bilnummerskilt) og **prompt-injection**. Grunnen til at prompt-injection er relevant: Rekognition-labels går inn i prompt-strengen uten sanering. En angriper kan skrive instruksjoner på et bilde, Rekognition fanger opp teksten som label, og instruksjonen havner i LLM-ens kontekst.

Policyen skal være knyttet til applikasjonens Bedrock-kall ved at `guardrailIdentifier` og `guardrailVersion` sendes med i invokasjonen — det holder ikke bare å ha den definert i Terraform.

Nyttige referanser:
- Terraform-ressursen: <https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/bedrock_guardrail>
- AWS Bedrock Guardrails (konsepter, policy-typer, eksempler): <https://docs.aws.amazon.com/bedrock/latest/userguide/guardrails.html>

Leveranser: Terraform-kode, to screenshots (ett blokkert og ett tillatt kall), og en begrunnelse på 200–300 ord i `README_SVAR.md` som forklarer hvilke policy-valg du gjorde og hvorfor de er fornuftige for *dette* scenariet.

### Oppgave 5 — Drøfteoppgave (15p)

En tradisjonell webapplikasjon er — med litt arbeid — deterministisk, testbar og enkel å regne på. En tjeneste som delegerer kjernefunksjonaliteten til en LLM, slik Wishifai gjør, er det ikke. Samme input gir ikke samme output, kvalitet er vanskelig å måle, og kostnaden varierer per request.

Dere skal skrive en refleksjon på **800–1000 ord** i `README_SVAR.md` om hva som er annerledes med DevOps-praksis for LLM-baserte tjenester sammenlignet med tradisjonelle webapplikasjoner. Mulige vinklinger: testing og evaluering av ikke-deterministisk oppførsel, prompt-versjonering som en del av CI/CD, cost- og token-observability, hvordan guardrails fungerer som en governance-mekanisme i leveransemodellen, hallusinering som et driftsproblem.

Bruk konkrete eksempler fra din egen besvarelse. En generell "DevOps er viktig"-essay gir ikke uttelling.



## Vurdering

| Oppgave | Poeng |
|---------|-------|
| 1 Container og CI | 20 |
| 2 ECS Express og IaC | 25 |
| 3 Observability | 20 |
| 4 Responsible AI | 20 |
| 5 Drøfteoppgave | 15 |
| **Totalt** | **100** |

Hvis du velger å ikke besvare en oppgave, nevn det i `README_SVAR.md`.
