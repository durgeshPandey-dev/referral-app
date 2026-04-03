const { execSync } = require("child_process");

const projects = JSON.parse(
  execSync("npx nx show projects --json", { encoding: "utf8" })
);

console.log("\nTargets:");
console.log(
  "  make generate                                     Auto generate swagger spec (go generate) and UI client (Kubb CLI)"
);
console.log(
  "  make run       PROJECT=<project> TARGET=<target>  Run a specific target on a specific project"
);
console.log(
  "  make all       TARGET=<target>                    Run a target across all projects"
);
console.log(
  "  make affected  TARGET=<target>                    Run a target only on projects affected by current git changes"
);
console.log(
  "  make ci                                           Run lint, typecheck, build and test across all projects\n"
);
console.log("Available projects and their targets:\n");

for (const project of projects) {
  const data = JSON.parse(
    execSync(`npx nx show project ${project} --json`, { encoding: "utf8" })
  );
  const targets = Object.keys(data.targets || {});
  console.log(`  ${project}`);
  targets.forEach((t) => console.log(`    - ${t}`));
  console.log("");
}
