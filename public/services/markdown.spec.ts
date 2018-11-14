import * as markdown from "./markdown";

const testCases = [
  {
    input: "Visit [GitHub](https://github.com) to learn more.",
    expectedFull: '<p>Visit <a href="https://github.com" target="_blank">GitHub</a> to learn more.</p>',
    expectedSimple: '<p>Visit <a href="https://github.com" target="_blank">GitHub</a> to learn more.</p>'
  },
  {
    input: "My Picture ![](http://demo.dev.fider.io:3000/images/100/28)",
    expectedFull: '<p>My Picture <img src="http://demo.dev.fider.io:3000/images/100/28" alt=""/></p>',
    expectedSimple: "<p>My Picture </p>"
  },
  {
    input: "# Hello World",
    expectedFull: "<h1>Hello World</h1>",
    expectedSimple: "<p>Hello World</p>"
  },
  {
    input: "~~Option 3~~",
    expectedFull: "<p><del>Option 3</del></p>",
    expectedSimple: "<p><del>Option 3</del></p>"
  },
  {
    input: "Check this out: `HEEEY`",
    expectedFull: "<p>Check this out: <code>HEEEY</code></p>",
    expectedSimple: "<p>Check this out: <code>HEEEY</code></p>"
  },
  {
    input: `# Hello World
How are you?`,
    expectedFull: `<h1>Hello World</h1>
<p>How are you?</p>`,
    expectedSimple: `<p>Hello World</p><p>How are you?</p>`
  }
];

testCases.forEach(x => {
  test(`Can parse markdown ${x.input} to ${x.expectedFull} (full mode)`, () => {
    const result = markdown.full(x.input);
    expect(result).toEqual(x.expectedFull);
  });
});

testCases.forEach(x => {
  test(`Can parse markdown ${x.input} to ${x.expectedSimple} (simple mode)`, () => {
    const result = markdown.simple(x.input);
    expect(result).toEqual(x.expectedSimple);
  });
});
