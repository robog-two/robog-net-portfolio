---
title: "Generative&nbsp;AI Ethical&nbsp;Framework"
layout: docs.njk
tags:
  - docroot
---

**This document is an outline of my personal best practices and beliefs when incorporating Generative AI into my work.**

In a stark contrast to [Satya Nadella of Microsoft](https://snscratchpad.com/posts/looking-ahead-2026/), I believe that there *is* a real problem with the prevalence of so-called "AI slop"[^1] on the internet today. I'm someone who has always advocated for user privacy, open source, and using technology as a means to enable human joy and well-being. It is clear to me that the current state of Generative AI usage across the Fortune 500 and the internet broadly is increasingly negative and problematic. In this document, I'd like to outline my own policy on using Generative AI, with several core axioms[^2] and reasoning behind each.

This isn't a legally binding document or contract, and my stance on the correct way to incorporate this emerging technology according to my own moral principles is evolving—sometimes I use it wrong, too. However, I think it's important to hold myself accountable, and if you notice older code that appears to be generated without disclosure or contains LLM-created fabrications, please [reach out to me personally](https://robog.net/about) and I will make an effort to correct it.

# 1. The usage of Generative AI should be *disclosed*.

I feel like this principle makes sense in the context of my other values. I believe that collection of private user data should be disclosed, I believe that the source code of software we use should be disclosed, and I believe that the practices and thoughts behind that development should be disclosed. Sharing information is part of creating an equitable world. When we are informed we have agency, we can choose where to put our time and effort, what we can trust, and who we can count on.

I expect that some people will see a Generative AI disclosure and immediately turn the other way. I think that's okay, in the same way that I think it's okay that some people choose not to use social media, or some people choose not to eat meat. Everyone has their own preferences and beliefs, and they should be allowed to live according to those standards. *They cannot make this choice for themselves when reading or interacting with content created by Generative AI that has not been disclosed*.

I also understand that there is some nuance to this, based on whether or not I have used Generative AI for my own personal reasons or have incorporated it into the final product. For example, I generally do not disclose Generative AI usage in the following circumstances, when:

- Generative AI is used for researching or personal learning—think, an AI bot that answers questions about documentation for a library—but the final product comes from my own original ideas
- Generative AI is used to advise or comment upon an existing work, without directly contributing to the work in any way
- The output  of Generative AI *is* incorporated into the work but only affects stylistic or formatting choices, such as a spelling or grammar checker, or a code formatter. I will only choose not to disclose this kind of usage if the results did not substantially affect the work, I have tracked the changes made, and I have verified and take accountability for the work produced. For example, I sometimes use the Mac "Proofread" feature, which may add or remove commas and correct spelling, but I always vet each individual edit, character by character.

When I do disclose, I am clear and specific about what models are used in what ways and to what extent.

# 2. The usage of Generative AI should be *unobtrusive*.

When I hear the words "[intelligence as a design material](https://www.threads.com/@zuck/post/DR0Wsy0EhCP/today-were-establishing-a-new-creative-studio-in-reality-labs-led-by-alan-dye?hl=en)," I immediately understand that this is not the words of a designer but the words of a salesperson. This is in fact the inverse of design as a process: When you treat the technology as the basis of design, you are not making decisions about how technology can fit effortlessly into our lives, you are putting technology into places it does not belong. Therefore, design should serve to hide technology—creating something that serves a user-oriented purpose—and must use Generative AI in a way that is nearly invisible. The user sees valuable, creative, truthful, and original ideas, and the technology stays in the background to facilitate that goal.

# 3. Generative AI should be a *tool*.

The same general ideas that guide other forms of craftsmanship should also guide the creation of software with agentic coding tools. These new tools are powerful, but they also require proper safety precautions. Just like a belt sander, they may sand off the corners if not applied gently and carefully. Like a drill, they can easily go too far and create cracks in drywall when they are not given an anchor.

In literal terms, agentic coding tools and the work they produce should be
1) Limited in the amount of information they are allowed to access,
2) Limited in the scope of what they are allowed to change, and
3) Carefully inspected and double-checked for security, completeness, and clarity.

In other words, if Generative AI is a tool, always *measure twice and cut once*.

# 4. The usage of Generative AI should respect the *environment* and the *individual*.

I don't want to under– or oversell the environmental impact of large language models. As of right now, the only thing that is clear about their environmental impact is that the scope and intensity of the impact is unclear. This is because the developers of foundation models do not comprehensively track the resources that they are using nor do they make the usage data they do track public. There is a *severe lack of transparency* from all major producers of AI technology in this regard. Another problem this creates is privacy: It's often unclear if images uploaded to an AI chat or API service may be used to train diffusion models in the future, or if chats containing sensitive health or legal information will be exposed to employees of the company.

There is a clear solution to both of these problems: Use smaller, more efficient models, and host them on personally owned infrastructure. The technology needed to achieve this is complex, because smaller models often require clever inclusion of context (i.e. Retrieval Augmented Generation), or very fine-tuned prompting. However, the advantages of privacy and efficiency are clearly worth the additional design effort.

## Final Note on AI "Art"
**I do not and will not use or publish AI generated images or video in my work. Period.**
When I need art assets and I cannot create them myself, I will often consult permissively licensed or copyleft media sources, like the Library of Congress or Wikimedia Commons, giving attribution as the license requires. While I don't agree with the way that people talk about AI generated art—it is mostly through a framework of copyright law, which in general tends to support the rights of large media conglomerates and not small creators—I still think it is important to recognize that the technology is both bad for the environment and has been developed quite unethically, and most proponents of AI "art" detest actual artists and their livelihoods. I sort of feel about AI art the same way I feel about Cryptocurrency: Cool math, not so cool intended purpose.[^3]

At the end of the day, I want to create things that are joyful and useful. Sometimes large language models support this goal, and sometimes they come into conflict with it. My hope is that the axioms outlined in this document make it clear that I put a lot of thought into how and when I choose to include this powerful new technology in my own projects.

I hope you find them joyful and useful.

Best,<br/>
Sam

<p style="text-align: center; margin-top: 50px; margin-bottom: 120px; width: 100%;">//////////// 🫀 ////////////</p>

<p style="text-align: center; margin-bottom: 350px; width: 100%;">Originally written on January 9th, 2026 after long deliberation and discussions with friends and those close to me.<br/>Thank you, Rob. You have a brilliant and wonderful mind.<br/><br/>This document was not created by Generative AI.<br/>Claude Sonnet 4.5 and Gemma 3 were consulted for feedback on the prose and structure of inital drafts, but all text was produced originally and finally by my own human fingers.<br/>I would know, I can feel my carpal tunnels.</p>


*Footnotes:*
[^1]: I have my own qualms with the term "AI Slop." I think it is applied very broadly and generally used to imply that Generative AI can only create useless things. I know for a fact that this isn't true, in other words, I think that "slop" is far more of a problem than "AI." Unfortunately, when the norm has been so polluted by [irresponsible](https://techcrunch.com/2025/06/12/the-meta-ai-app-is-a-privacy-disaster/) and [dangerous](https://www.nytimes.com/2025/09/11/technology/google-meta-chatgpt-ai-chatbots.html) usages of this technology, it's understandable that it has a negative reputation.

[^2]: These axioms were partially inspired by Dieter Rams' 10 Principles of Good Design. I think his work is quite thought-provoking and I find it valuable for the way in which it integrates humanistic principles with minimalism, and starts conversations about maximizing benefit to people while minimizing waste and resource usage. Also, I have them hung up on my wall, so I kind of think about them a lot, no matter what I'm writing.

[^3]: I know this wasn't a very specific example of the misuses of those technologies, so let me elaborate. The "cool math" is diffusion models or zero-trust cryptography. The "not so cool purpose" is pornography, counterfeiting, money laundering, and trafficking. It's easy to separate the ethical uses of LLMs from the unethical ones. They are good at cheating, lying, and plagiarizing, but I think this is outweighed by their capacity to reason, code, communicate, and express—but not feel—emotions. Plus, there is a clear and concerted effort in the LLM space, especially by Anthropic, to reduce these specific harms in a way that I have not seen with image generation models. [Grok/xAI's behavior and policies on this matter are dangerous and frankly disgusting.](https://www.theverge.com/news/859715/x-grok-ai-deepfakes)
