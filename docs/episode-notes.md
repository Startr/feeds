---
title: "Run It Local - Episode Notes"
date: 2026-04-10
status: DRAFT
hosts: Izzy Plante (CEO), Alex Somma (CTO)
show: Run It Local
---

# Startr.it - Episode Notes

## Show overview

**Show:** Startr.it
**Hosts:** Izzy Plante (Co-founder & CEO) and Alex Somma (Co-founder & CTO), Sage
**Format:** Two-host conversation (pilot format). Roundtable with guests from ep 3+.
**Target runtime:** 30-35 minutes per episode
**Thesis:** Real autonomy vs. compliance theater
**Voice:** Direct, operator-flavored, assumed-intelligent. Same voice for engineers
and parents. No code-switching, no dumbing down.

---

## Episode 1: "The Licence Nobody Teaches"

**Source:** https://sage.is/resources/the-licence-nobody-teaches/
**Working title:** "The Licence Nobody Teaches"
**One-liner:** The MIT licence is 170 words. It gives everything away. The AGPL
protects everything. Most CS graduates have never heard of it.

### Why this is Episode 1

This episode is meta in the best way. Startr/feeds (the tool publishing this
podcast) is AGPL-3.0. The show thesis is "real autonomy vs. compliance theater."
Licensing IS the mechanism. If you want to talk about who owns your data, who
controls your infrastructure, and who profits from your work, the licence is
where that fight actually happens. Start here. Everything else follows.

Also: Alex is directly quoted in the article. This isn't commentary on someone
else's work. This is the show explaining its own foundation.

### The story in 30 seconds

In 2013, GitHub's co-founder stood on stage and called the GPL "not freedom"
because it was long. GitHub then built choosealicense.com with MIT as the default.
In a single decade, open-source licensing inverted from 59% copyleft to 40%
permissive. The result: Amazon took Redis, Elastic, and Terraform, built
billion-dollar services on them, returned nothing, and the original creators
had no recourse. Every one of them eventually switched to AGPL or similar.
The licence they should have chosen from day one.

### Episode structure

| Section                 | Time        | What happens                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| ----------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Cold open               | 0:00-1:30   | Hook: "The MIT licence is 170 words. It fits on a business card. And it's the reason Amazon runs your code as a billion-dollar service without paying you a cent. Most developers chose it without reading it. Most CS programs never taught the alternative. We're going to talk about why."                                                                                                                                                                                                       |
| Intro                   | 1:30-3:30   | Who Izzy and Alex are. What Startr.Media and Sage is (one sentence). What this episode is about: commentary on the gap between what companies call open source permissive and what it actually isn't - protective.                                                                                                                                                                                                                                                                                  |
| The GitHub origin story | 3:30-8:00   | GitHub co-founder, Preston-Werner's 2013 OSCON keynote: "MIT is Nice and short" vs "AGPL3 is long and that's not what freedom is. (14,000 words)" choosealicense.com putting MIT first. The consequence was a noticeable increase in permissive licenses and an decrease in protective ones. This was not organic. It was designed. Alex's quote: "The curriculum doesn't teach licensing because the companies that fund the curriculum benefit from the licence the students choose by default."  |
| The curriculum gap      | 8:00-12:00  | No major CS program teaches comparative licensing. Students graduate having never read a licence. Google, Amazon, Meta, Microsoft fund the programs, sponsor the hackathons, ship the tools (React, TensorFlow, VS Code, PyTorch), all permissive. The ecosystem normalizes MIT as default. The students carry the default with them into their professional lives.                                                                                                                                 |
| The Redis story         | 12:00-18:00 | Salvatore Sanfilippo creator of Redis. chose a BSD licence. Amazon builds ElastiCache. Returns nothing. Redis relicenses to SSPL in 2024. Community fractures. Valkey fork backed by Amazon. External contributors with 5+ commits drops to zero. Sanfilippo returns December 2024. Redis adds AGPL May 2025. Sixteen years to learn one lesson. **"The licence you choose on day one determines who profits on day ten thousand."**                                                                |
| The pattern             | 18:00-22:00 | It's not just Redis. Walk through: Elastic (Apache 2.0 -> Amazon fork -> SSPL -> AGPL). HashiCorp (MPL -> BSL -> OpenTofu fork). Same arc every time: permissive -> adoption -> extraction -> panic -> relicense -> community fracture. Every case preventable with AGPL from day one.                                                                                                                                                                                                              |
| The Google ban          | 22:00-26:00 | Google banned AGPL in 2010 for reasons specific to Google's monorepo architecture. Rational for them. Then the industry copied it without understanding why. Heather Meeker: "The single most significant misconception about copyleft licenses is caused by the use of the word 'viral.'" The metaphor is wrong. Copyleft does not cross process boundaries. "One company banned AGPL because it worked. The industry banned AGPL because that company did. Nobody in the chain read the licence." |
| The comeback of AGPL    | 26:00-29:00 | RedMonk 2026: permissive dropped from 82% to 73%. First reversal in a decade. Elastic -> AGPL. Redis -> AGPL. Nextcloud 250M euros AGPL infrastructure. Vitalik Buterin's public reversal: "Nonzero openness is the only way the world does not eventually converge to one actor controlling everything." The AI dimension: MIT code trains proprietary models without attribution.                                                                                                                 |
| Why we chose AGPL       | 29:00-32:00 | Bring it home. Startr/feeds is AGPL-3.0. Sage.is AI-UI is AGPL-3.0. This show is published by an AGPL tool. Why: if anyone takes this code, runs it as a service, the modifications are shared back. The embed-vs-fork distinction (Plausible model). We have created our "alternative 170 words" and a comparison table CS educators can use to explain : MIT vs GPL vs AGPL. All three allow commercial use. The only difference is whether improvements come back.                               |
| Close                   | 32:00-35:00 | The full article is on (sage.is/resources). Next episode preview: "From Scratch" -- what happened when we removed the workbench from the classroom and what Sage.Education AI is doing to bring it back. Subscribe to our RSS feed. Specific ask: "Next time someone tells you to use MIT because it's simpler, ask them: simpler for whom?"                                                                                                                                                        |

### Key quotes to hit

- "Nice and short" / "That's not what freedom is." (Preston-Werner, 2013)
- "The curriculum doesn't teach licensing because the companies that fund the
  curriculum benefit from the licence the students choose by default." (Alex Somma)
- "The licence you choose on day one determines who profits on day ten thousand."
- "The single most significant misconception about copyleft licenses is caused by
  the use of the word 'viral'." (Heather Meeker)
- "One company banned AGPL because it worked. The industry banned AGPL because
  that company did. Developers avoid AGPL because the industry does. Nobody in
  the chain read the licence."
- "Nonzero openness is the only way the world does not eventually converge to one
  actor controlling everything." (Vitalik Buterin)

### The comparison table (show notes + say on air)

|  | MIT | GPL v3 | AGPL v3 |
|---|---|---|---|
| Commercial use | Yes | Yes | Yes |
| Modification | Yes | Yes | Yes |
| Distribution | Yes | Yes | Yes |
| Must share modifications | No | Yes (if distributed) | Yes (if distributed OR served) |
| SaaS loophole closed | No | No | Yes |
| Protects against cloud extraction | No | Partially | Yes |

### Show notes (publish with episode)

**Episode 1: The Licence Nobody Teaches**

The MIT licence is 170 words. It gives everything away. The AGPL protects
everything. Most CS graduates have never heard of it.

Izzy and Alex walk through the decade-long inversion of open-source licensing,
the Redis/Elastic/Terraform pattern, why Google's AGPL ban became an industry
religion, and why the copyleft comeback matters. Plus: why Startr/feeds and
Sage.is are AGPL-3.0, and what that means for anyone who wants to self-host.

**Stories covered:**
- 0:00 -- The 170 words that changed open source
- 3:30 -- How GitHub made MIT the default
- 8:00 -- The curriculum gap nobody talks about
- 12:00 -- What Amazon did to Redis (and Elastic, and Terraform)
- 22:00 -- Google's AGPL ban and the viral myth
- 26:00 -- The copyleft comeback (2025-2026)
- 29:00 -- Why we chose AGPL for everything

**Source article:** [The Licence Nobody Teaches](https://sage.is/resources/the-licence-nobody-teaches/)

**Links mentioned:**
- [choosealicense.com](https://choosealicense.com) (the origin story)
- [RedMonk 2026 State of Open Source Licensing](https://redmonk.com)
- [Vitalik Buterin on copyleft](https://vitalik.eth.limo)
- [Heather Meeker, Open Source for Business](https://heathermeeker.com)
- [Startr/feeds on GitHub](https://github.com/Startr/feeds) (AGPL-3.0)
- [Sage.is](https://sage.is)

### Pre-record checklist

- [ ] Both hosts have re-read "The Licence Nobody Teaches" fresh
- [ ] Bullet-point outline in a shared doc (not a script)
- [ ] Test recording (5 min) to confirm levels and vibe
- [ ] Know the Redis timeline cold: 2009 BSD, 2024 SSPL, 2025 AGPL
- [ ] Know the choosealicense.com origin: 2013 Preston-Werner keynote
- [ ] Know the numbers: 59%/41% -> 40%/10% inversion. 82% -> 73% reversal.
- [ ] Have the comparison table visible during recording

---

## Episode 2: "From Scratch"

**Source:** https://sage.education/posts/blog/en/from-scratch/
**Working title:** "From Scratch"
**One-liner:** We removed the workbench from the classroom. AI is finishing
the job. The alternative is not to reject AI. It is to refuse the vending
machine architecture.

### Why this follows Episode 1

Episode 1 is about who owns the code. Episode 2 is about who uses it, and how.
The licensing episode establishes the infrastructure thesis (own your tools).
This episode extends it to education (own your learning). Same underlying
concern: the difference between tools that make you capable and systems that
make you dependent. AGPL protects the code. Maker-oriented AI protects the
learner. Both are about real autonomy.

Also: this hits the education/parent audience directly. Episode 1 is more
operator-sided (licensing, infrastructure). Episode 2 balances the other way
(education, children, the homeschool parent listening). Same voice, different
entry point. Together they establish the show's range.

### The story in 30 seconds

Over 30 years, schools removed shop class, home economics, and hands-on making.
Students learned to choose, not to make. Now AI is the universal cognitive
offloader: insert a prompt, receive a finished product. 97% of teen screen time
is consumption. AI perfects that ratio. But the maker-oriented alternative
exists. It looks like giving people tools, not answers. A workbench, not a
vending machine.

### Episode structure

| Section | Time | What happens |
|---|---|---|
| Cold open | 0:00-1:30 | The cookie story. Baking from scratch with a fourteen-year-old. Eggshell in the bowl. "Why can't we just buy cookies?" The cookie was not the point. The process was the point. The mess was the point. |
| Connect to Episode 1 | 1:30-3:00 | "Last episode we talked about who owns the code. Today we're talking about who uses it, and how. Same question, different lens. The licence protects the tool. The pedagogy protects the learner." |
| The disappeared classroom | 3:00-8:00 | 1982: most high schools had woodworking, metalworking, auto shop, home ec. 2013: fewer than a quarter. What happened: No Child Left Behind tied funding to test scores. Race to the Top. Districts cut non-tested subjects. Workshops became computer labs. Matthew Crawford: "The disappearance of tools from our common education is the first step toward a wider ignorance -- not just of how to fix things, but of the very idea that things can be fixed." |
| The pre-made generation | 8:00-12:00 | Food prep dropped from 2+ hours/day to under 1 hour. Only 10% of millennials bake from scratch. Apple's planned obsolescence. "The pre-made cookie and the pre-made answer have the same architecture: a finished product that arrives without requiring the consumer to understand how it was produced." |
| AI as universal vending machine | 12:00-19:00 | ChatGPT: 100M users in 2 months. 1 in 4 American teens using it for schoolwork. 11% of submissions flagged by Turnitin. But the real problem isn't cheating. Ethan Mollick (Wharton): students who used AI as a crutch scored weaker on follow-up assessments. Anthropic's RCT: 52 junior engineers, AI-reliant group scored 17% lower on comprehension. "The tool improved the assignment. It degraded the learning." Cognitive offloading (Sparrow 2011, Risko 2016). Only 3% of teen screen time is creation. 97% is consumption. AI does not correct this ratio. It perfects it. |
| The workbench that scales | 19:00-25:00 | The constructive alternative. Mitch Resnick (MIT): digital tools as instruments of creation, not consumption. The trajectory from Scratch to Snap!. Paulo Blikstein (FabLearn Labs): maker-oriented programs improve problem-solving, especially for underrepresented students. Kylie Peppler: "When the measure of success is 'did you build something that works' rather than 'did you select the correct answer,' a different population of students succeeds." The central question: does the tool hand over raw materials and say "build," or a finished product and say "approve"? |
| Sage's position | 25:00-29:00 | Connect to Sage.is AI-UI. Open source (AGPL, from last episode). Self-hostable. Model-agnostic. Workshops for building agents/tools/knowledge bases. "A vending machine teaches you to insert coins. A workshop teaches you to use tools." Gary Stager: does the student use AI like a woodworker uses a lathe (extending capability while requiring skill) or like a customer uses a catalogue (selecting finished products)? |
| The education stakeholder lens | 29:00-32:00 | Speak directly to the parent/educator audience. Homeschool parents, microschool operators, public school teachers. The shared concern: "approved" and "compliant" tools are not the same as safe tools. A seven-year-old cannot consent to having their writing assignments sent to Meta as training data. Learning data is uniquely sensitive. The CLOUD Act doesn't stop at adults. Connect back to Episode 1: the licence matters because the tool matters because the child using the tool matters. |
| Close | 32:00-35:00 | Return to the cookie. "The empty workshop is not waiting for better tools. It is waiting for someone to pick the tools up." Where to find the article (sage.education). Subscribe. Specific ask: "If you're a parent, a teacher, or a builder, the next time someone offers you a finished answer, ask: what did I learn by receiving it? And what would I have learned by building it myself?" Next episode preview. |

### Key quotes to hit

- "The cookie was not the point. The process was the point. The mess was the point."
- "The disappearance of tools from our common education is the first step toward a
  wider ignorance -- not just of how to fix things, but of the very idea that things
  can be fixed." (Matthew Crawford)
- "The pre-made cookie and the pre-made answer have the same architecture."
- "The tool improved the assignment. It degraded the learning." (Ethan Mollick)
- "Only 3% is content creation. The rest is consumption."
- "AI is the universal cognitive offloader. It does not offload one skill.
  It offloads all of them."
- "When the measure of success is 'did you build something that works' rather than
  'did you select the correct answer,' a different population of students succeeds."
  (Kylie Peppler)
- "The alternative is not to reject AI. It is to refuse the vending machine
  architecture."
- "The empty workshop is not waiting for better tools. It is waiting for someone
  to pick the tools up."

### Research numbers to know

- 1982: majority of US high schools had shop programs. 2013: fewer than 25%.
- AP exam volume: ~1M (2000) to 5M+ (2020)
- Food prep time: 2+ hours/day (1965) to under 1 hour (2000s)
- Only 10% of millennials frequently bake from scratch
- ChatGPT: 100M users in 2 months (Nov 2022)
- 1 in 4 American teens use AI for schoolwork (Pew 2024)
- 11% of submissions flagged by Turnitin (mid-2023)
- Anthropic RCT: 17% lower comprehension scores for AI-reliant group
- Teen screen time: 8h39m/day, only 3% creation
- EU: extending device lifespan by 1 year saves ~4M tonnes CO2/year by 2030

### Show notes (publish with episode)

**Episode 2: From Scratch**

We removed shop class, home economics, and hands-on making from schools over
thirty years. We created a generation of consumers, not creators. Now AI is
repeating the same pattern at industrial speed. The alternative exists. It
looks like giving people tools, not answers.

**Stories covered:**
- 0:00 -- The cookie (and why the mess is the point)
- 3:00 -- The disappeared classroom: what happened to shop class
- 8:00 -- The pre-made generation: from making to consuming
- 12:00 -- AI as universal vending machine (and the research that proves it)
- 19:00 -- The workbench that scales: maker-oriented alternatives
- 25:00 -- What Sage is building (and why it's AGPL)
- 29:00 -- A message for parents and teachers

**Source article:** [From Scratch](https://sage.education/posts/blog/en/from-scratch/)

**Research cited:**
- Anthropic RCT on AI-assisted coding (January 2025)
- Mollick, Wharton AI + education research (2023)
- Sparrow et al., "Google Effects on Memory" (Science, 2011)
- Risko & Gilbert, "Cognitive offloading" (Trends in Cognitive Sciences, 2016)
- Common Sense Media, teen screen time report (2021)
- Crawford, *Shop Class as Soulcraft* (2009)
- Resnick, *Lifelong Kindergarten* (2017)

**Links mentioned:**
- [Sage Education](https://sage.education)
- [Sage.is AI-UI](https://sage.is)
- [Startr/feeds on GitHub](https://github.com/Startr/feeds)

### Pre-record checklist

- [ ] Both hosts have re-read "From Scratch" fresh
- [ ] Bullet-point outline in a shared doc (not a script)
- [ ] Test recording (5 min) to confirm levels
- [ ] Know the key numbers: shop class decline, teen screen time, Anthropic 17%
- [ ] Have the cookie story beats clear (don't over-rehearse, but know the arc)
- [ ] Decide: does Izzy or Alex tell the cookie story? (Whoever it happened to)

---

## Production notes (both episodes)

**Recording:** One continuous take if possible. Edit only for dead air, not for pacing.
**Audio:** 128kbps minimum, mono fine. USB condenser mic or better.
**Transcript:** Whisper (runs locally, free). Review for proper nouns: "Perplexity,"
"Redis," "Sanfilippo," "Buterin," "Peppler," "Blikstein," "Resnick."
**Cover art:** Placeholder for launch; iterate later.
**Upload:** Spotify for Podcasters. Wait for auto-feed. Run the rewriter. Verify.
**Submit feed:** Apple Podcasts Connect + Spotify for Podcasters.
**Embed:** Player on sage.is/podcast (readiest site for Episode 1 since it sources
from sage.is/resources). sage.education/podcast for Episode 2.

### Episode pairing logic

These two episodes are designed as a pair. Episode 1 (licensing) is the
infrastructure thesis: who owns the tools. Episode 2 (education) is the human
thesis: who uses the tools, and how. Together they establish the show's range
for both audiences (operators and parents) without code-switching.

Record them in the same session if energy allows. Two 30-minute takes with a
break between. Ship Episode 1 immediately, Episode 2 one week later. Weekly
cadence starts from day one.

### Updated episode order (vs. design doc)

The design doc at `docs/podcast-site-design.md` originally planned Episode 1 as
"The Prompt You Thought Was Private" (Perplexity story). That episode becomes
Episode 3 or later. Licensing is the stronger foundation: it sets up the show's
thesis mechanically, and it's meta (the show's own tool is AGPL-3.0).

| Episode | Title | Source |
|---|---|---|
| 1 | The Licence Nobody Teaches | sage.is/resources/the-licence-nobody-teaches |
| 2 | From Scratch | sage.education/posts/blog/en/from-scratch |
| 3+ | The Prompt You Thought Was Private | sage.is/resources (original Ep 1 plan) |
| 3+ | The $70,000 Illusion | sage.is/resources (original supporting story) |
| 3+ | The Thinking Layer | sage.is/resources (original supporting story) |
