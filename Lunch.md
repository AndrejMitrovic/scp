# A round of lunch

This repo includes a demonstration program called `lunch`.
It simulates a group of friends or coworkers
(organized into some configurable network of quorum slices)
coming to consensus on what to order for lunch.

This document contains the
(simplified)
output of one round of lunch consensus.
It describes step by step how the network goes from conflicting nominations to agreeing on an outcome.

This example was generated using the [“3 tiers”](https://github.com/bobg/scp/blob/master/cmd/lunch/toml/3tiers.toml) network configuration,
in which:
- a top tier of nodes — Top1, Top2, Top3, and Top4 — each depends on two of its neighbors for a quorum;
- a middle tier — Mid1, Mid2, Mid3, and Mid4 — each depends on any two members of the top tier; and
- a bottom tier — Low1 and Low2 — each depends on any two members of the middle tier.

Note that lines ending with `-> ∅`
(meaning the node produces no output in response to a protocol message)
are ones the `lunch` program does not normally display,
even though some of them do update the network’s internal state.
They are included here to make the workings of the consensus algorithm clearer.

```
Top4: ∅ -> (Top4 NOM X=[salads], Y=[])
```

Top4 votes to nominate salads.
In an SCP protocol message,
X is the set of nominees voted for.

```
Mid1: ∅ -> (Mid1 NOM X=[burgers], Y=[])
```

Mid1 votes to nominate burgers.

```
Mid3: ∅ -> ∅
```

Mid3 would like to nominate something,
but it’s not in its own “high-priority neighbors” list
(at this particular time for this particular slot),
so it discards its own nomination.

```
Top2: ∅ -> (Top2 NOM X=[indian], Y=[])
Top1: ∅ -> (Top1 NOM X=[burritos], Y=[])
```

Top2 and Top1 vote to nominate Indian food and burritos, respectively.

```
Top3: ∅ -> ∅
```

Top3, like Mid3, is not in its own high-priority-neighbors list,
so it does not have the ability to nominate anything at the moment.

```
Top3: (Top4 NOM X=[salads], Y=[]) -> (Top3 NOM X=[salads], Y=[])
```

However,
Top4 _is_ one of Top3’s high-priority neighbors.
It sees its vote to nominate salads and it echoes it.

```
Top4: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Mid1: (Top4 NOM X=[salads], Y=[]) -> ∅
Top2: (Top4 NOM X=[salads], Y=[]) -> ∅
Top4: (Top2 NOM X=[indian], Y=[]) -> ∅
Top3: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top1: (Top4 NOM X=[salads], Y=[]) -> ∅
```

These nodes all see nominations from their peers but,
throttled by the priority mechanism,
do not echo them.

```
Mid3: (Top4 NOM X=[salads], Y=[]) -> (Mid3 NOM X=[salads], Y=[])
```

Mid3 does echo Top4’s nomination.

```
Mid1: (Top2 NOM X=[indian], Y=[]) -> ∅
Top3: (Top2 NOM X=[indian], Y=[]) -> ∅
Top2: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top1: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Mid3: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top4: (Top1 NOM X=[burritos], Y=[]) -> ∅
Mid1: (Top1 NOM X=[burritos], Y=[]) -> ∅
Top4: (Top3 NOM X=[salads], Y=[]) -> ∅
```

Top4 sees Top3’s nomination of salads.
Top3 may or may not be one of Top4’s high-priority neighbors at the moment,
but in any case Top4 has already nominated salads himself,
so this message from Top3 does not cause any change in Top4’s state
(which means Top4 sends out no new message in response).

```
Top3: (Top1 NOM X=[burritos], Y=[]) -> ∅
Mid3: (Top2 NOM X=[indian], Y=[]) -> ∅
Top2: (Top1 NOM X=[burritos], Y=[]) -> ∅
Mid1: (Top3 NOM X=[salads], Y=[]) -> ∅
Top4: (Mid3 NOM X=[salads], Y=[]) -> ∅
Top1: (Top2 NOM X=[indian], Y=[]) -> ∅
Top3: (Mid3 NOM X=[salads], Y=[]) -> ∅
Top2: (Top3 NOM X=[salads], Y=[]) -> ∅
Mid1: (Mid3 NOM X=[salads], Y=[]) -> ∅
Mid3: (Top1 NOM X=[burritos], Y=[]) -> ∅
Top1: (Top3 NOM X=[salads], Y=[]) -> ∅
Top2: (Mid3 NOM X=[salads], Y=[]) -> ∅
Mid3: (Top3 NOM X=[salads], Y=[]) -> ∅
Top1: (Mid3 NOM X=[salads], Y=[]) -> ∅
```

More throttled-or-redundant peer nominations.

```
Mid4: ∅ -> (Mid4 NOM X=[salads], Y=[])
```

Mid4 votes to nominate salads.

```
Top1: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid4: (Top4 NOM X=[salads], Y=[]) -> ∅
Mid4: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top3: (Mid4 NOM X=[salads], Y=[]) -> ∅
Top4: (Mid4 NOM X=[salads], Y=[]) -> ∅
Top2: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid1: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid3: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid4: (Top2 NOM X=[indian], Y=[]) -> ∅
Mid4: (Top1 NOM X=[burritos], Y=[]) -> ∅
Mid2: ∅ -> ∅
Mid4: (Top3 NOM X=[salads], Y=[]) -> ∅
Mid4: (Mid3 NOM X=[salads], Y=[]) -> ∅
Mid2: (Top4 NOM X=[salads], Y=[]) -> (Mid2 NOM X=[salads], Y=[])
Mid3: (Mid2 NOM X=[salads], Y=[]) -> ∅
Mid2: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Mid1: (Mid2 NOM X=[salads], Y=[]) -> ∅
Top4: (Mid2 NOM X=[salads], Y=[]) -> ∅
Mid4: (Mid2 NOM X=[salads], Y=[]) -> ∅
Top1: (Mid2 NOM X=[salads], Y=[]) -> ∅
Top3: (Mid2 NOM X=[salads], Y=[]) -> ∅
Top2: (Mid2 NOM X=[salads], Y=[]) -> ∅
Mid2: (Top2 NOM X=[indian], Y=[]) -> ∅
Mid2: (Top1 NOM X=[burritos], Y=[]) -> ∅
Mid2: (Top3 NOM X=[salads], Y=[]) -> ∅
Mid2: (Mid3 NOM X=[salads], Y=[]) -> ∅
Mid2: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid1: (Low1 NOM X=[pizza], Y=[]) -> ∅
Top4: (Low1 NOM X=[pizza], Y=[]) -> ∅
Mid4: (Low1 NOM X=[pizza], Y=[]) -> ∅
Low1: (Top4 NOM X=[salads], Y=[]) -> ∅
Low1: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top2: (Low1 NOM X=[pizza], Y=[]) -> ∅
Mid3: (Low1 NOM X=[pizza], Y=[]) -> ∅
Top3: (Low1 NOM X=[pizza], Y=[]) -> ∅
Mid2: (Low1 NOM X=[pizza], Y=[]) -> ∅
Top1: (Low1 NOM X=[pizza], Y=[]) -> ∅
Low1: (Top2 NOM X=[indian], Y=[]) -> ∅
Low1: (Top1 NOM X=[burritos], Y=[]) -> ∅
Low1: (Top3 NOM X=[salads], Y=[]) -> ∅
Low1: (Mid3 NOM X=[salads], Y=[]) -> ∅
Low1: (Mid4 NOM X=[salads], Y=[]) -> ∅
Low1: (Mid2 NOM X=[salads], Y=[]) -> ∅
Low2: ∅ -> (Low2 NOM X=[sandwiches], Y=[])
Mid4: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top2: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Low1: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top4: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Low2: (Top4 NOM X=[salads], Y=[]) -> ∅
Mid2: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top3: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Mid1: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Mid3: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top1: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Low2: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Low2: (Top2 NOM X=[indian], Y=[]) -> ∅
Low2: (Top1 NOM X=[burritos], Y=[]) -> ∅
Mid1: (Top1 NOM X=[burritos], Y=[]) -> ∅
```

The nomination process continues,
with some nominations echoed,
others throttled,
and some nodes self-censoring.

```
Mid3: ∅ -> (Mid3 NOM X=[pasta salads], Y=[])
```

Mid3,
who originally wanted to nominate something but self-censored because it wasn’t in its own high-priority-neighbors list,
now is.
This happened because enough time elapsed since Mid3 began this nomination round that its high-priority list expanded.
It was previously echoing a nomination for pasta but now adds its own nomination, for salads.

```
Top2: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top3: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top2: (Top4 NOM X=[salads], Y=[]) -> ∅
Top3: (Mid3 NOM X=[salads], Y=[]) -> ∅
```

As mentioned,
Mid3 is now nominating both pasta and salads,
but an older protocol message of its,
from when it was still nominating only salads,
is only just now reaching Top3.
(Top3 does not respond because it is already nominating salads.)

```
Top1: ∅ -> ∅
Top3: (Mid2 NOM X=[salads], Y=[]) -> ∅
Top3: ∅ -> ∅
Mid3: (Top4 NOM X=[salads], Y=[]) -> ∅
Top3: (Top4 NOM X=[salads], Y=[]) -> ∅
Mid3: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top3: (Top2 NOM X=[indian], Y=[]) -> (Top3 NOM X=[indian salads], Y=[])
Mid1: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid3: (Top1 NOM X=[burritos], Y=[]) -> ∅
Top3: (Top1 NOM X=[burritos], Y=[]) -> ∅
Mid1: ∅ -> ∅
Mid3: (Top3 NOM X=[salads], Y=[]) -> ∅
Top1: (Top2 NOM X=[indian], Y=[]) -> (Top1 NOM X=[burritos indian], Y=[])
```

Top1’s high-priority-neighbor list has now expanded to include Top2,
and so it starts echoing its nomination of Indian food.

```
Mid3: (Top2 NOM X=[indian], Y=[]) -> ∅
Top3: (Mid4 NOM X=[salads], Y=[]) -> ∅
Top3: (Low1 NOM X=[pizza], Y=[]) -> ∅
Top3: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top2: (Top1 NOM X=[burritos], Y=[]) -> ∅
Top2: (Top3 NOM X=[salads], Y=[]) -> ∅
Top2: (Low1 NOM X=[pizza], Y=[]) -> ∅
Mid1: (Top4 NOM X=[salads], Y=[]) -> ∅
Top2: (Mid2 NOM X=[salads], Y=[]) -> ∅
Top2: ∅ -> ∅
Top2: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top2: (Mid3 NOM X=[salads], Y=[]) -> ∅
Top2: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid1: (Top2 NOM X=[indian], Y=[]) -> ∅
Mid1: (Low1 NOM X=[pizza], Y=[]) -> ∅
Mid1: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top1: (Mid3 NOM X=[salads], Y=[]) -> ∅
Mid1: (Top3 NOM X=[salads], Y=[]) -> ∅
Top1: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid1: (Mid3 NOM X=[salads], Y=[]) -> ∅
Top1: (Mid2 NOM X=[salads], Y=[]) -> ∅
Mid3: (Mid4 NOM X=[salads], Y=[]) -> ∅
Mid1: (Mid2 NOM X=[salads], Y=[]) -> ∅
Top1: (Top4 NOM X=[salads], Y=[]) -> ∅
Mid3: (Mid2 NOM X=[salads], Y=[]) -> ∅
Top1: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top4: (Mid2 NOM X=[salads], Y=[]) -> ∅
Mid3: (Low1 NOM X=[pizza], Y=[]) -> ∅
Top1: (Top3 NOM X=[salads], Y=[]) -> ∅
Top4: (Low1 NOM X=[pizza], Y=[]) -> ∅
Mid3: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top1: (Low1 NOM X=[pizza], Y=[]) -> ∅
Top4: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Top1: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Low1: ∅ -> (Low1 NOM X=[pizza], Y=[])
Top4: (Top2 NOM X=[indian], Y=[]) -> (Top4 NOM X=[indian salads], Y=[])
Top4: (Top1 NOM X=[burritos], Y=[]) -> ∅
Low2: (Top3 NOM X=[salads], Y=[]) -> ∅
Top4: (Mid3 NOM X=[salads], Y=[]) -> ∅
Top4: ∅ -> ∅
Top4: (Top3 NOM X=[salads], Y=[]) -> ∅
Top4: (Mid4 NOM X=[salads], Y=[]) -> ∅
Top4: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Top2: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Low2: (Mid3 NOM X=[salads], Y=[]) -> ∅
Top4: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Low1: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Low1: (Top3 NOM X=[indian salads], Y=[]) -> ∅
Top2: (Top3 NOM X=[indian salads], Y=[]) -> ∅
Mid1: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Mid4: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Top3: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Mid1: (Top3 NOM X=[indian salads], Y=[]) -> ∅
```

Plenty more of the same.

```
Top4: (Top3 NOM X=[indian salads], Y=[]) -> (Top4 NOM X=[salads], Y=[indian])
```

Something new: Top4 has moved Indian food from X,
the set of values he’s voting to nominate,
to Y, the set of values he _accepts_ as nominated.
This happens when either:

1. A _quorum_ votes-or-accepts the same value;
2. A _blocking set_ accepts it.

Top4 previously had seen Top2 vote to nominate Indian food and echoed that nomination.
Now that Top4 sees Top3 also voting to nominate Indian food,
condition 1 is satisfied:
Top2,
Top3,
and Top4 together form one of Top4’s quorums,
all voting for the same thing.

```
Mid3: (Top3 NOM X=[indian salads], Y=[]) -> ∅
Low2: (Mid4 NOM X=[salads], Y=[]) -> ∅
Low1: (Top1 NOM X=[burritos indian], Y=[]) -> ∅
Low2: (Mid2 NOM X=[salads], Y=[]) -> ∅
Low2: (Low1 NOM X=[pizza], Y=[]) -> ∅
Mid2: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Top1: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Mid4: (Top3 NOM X=[indian salads], Y=[]) -> ∅
Mid1: (Top1 NOM X=[burritos indian], Y=[]) -> ∅
Low1: (Top4 NOM X=[indian salads], Y=[]) -> ∅
```

More redundant-or-throttled nominations.

```
Top2: (Top1 NOM X=[burritos indian], Y=[]) -> (Top2 NOM X=[], Y=[indian])
Top3: (Top1 NOM X=[burritos indian], Y=[]) -> (Top3 NOM X=[salads], Y=[indian])
```

More “accepting” of the Indian-food nomination.

```
Mid1: (Top4 NOM X=[indian salads], Y=[]) -> ∅
Top4: (Top1 NOM X=[burritos indian], Y=[]) -> ∅
Mid4: (Top1 NOM X=[burritos indian], Y=[]) -> ∅
Mid3: (Top1 NOM X=[burritos indian], Y=[]) -> ∅
Low1: (Top4 NOM X=[salads], Y=[indian]) -> ∅
Top2: (Top4 NOM X=[indian salads], Y=[]) -> ∅
Mid4: (Top4 NOM X=[indian salads], Y=[]) -> ∅
Low2: (Mid3 NOM X=[pasta salads], Y=[]) -> ∅
Mid2: (Top3 NOM X=[indian salads], Y=[]) -> ∅
Mid2: (Top1 NOM X=[burritos indian], Y=[]) -> ∅
Top1: (Top3 NOM X=[indian salads], Y=[]) -> (Top1 NOM X=[burritos], Y=[indian])
Top1: (Top4 NOM X=[indian salads], Y=[]) -> ∅
Top3: (Top4 NOM X=[indian salads], Y=[]) -> ∅
```

Nomination continues.

```
Top4: (Top2 NOM X=[], Y=[indian]) -> ∅
```

Top4 sees that Top2 now “accepts” Indian food as nominated.
This will be important in a moment.

```
Top1: (Top4 NOM X=[salads], Y=[indian]) -> ∅
Top2: (Top4 NOM X=[salads], Y=[indian]) -> ∅
Mid1: (Top4 NOM X=[salads], Y=[indian]) -> ∅
Mid2: (Top4 NOM X=[indian salads], Y=[]) -> (Mid2 NOM X=[salads], Y=[indian])
Mid4: (Top4 NOM X=[salads], Y=[indian]) -> ∅
Low2: (Top3 NOM X=[indian salads], Y=[]) -> ∅
Mid1: (Top2 NOM X=[], Y=[indian]) -> ∅
Mid3: (Top4 NOM X=[indian salads], Y=[]) -> (Mid3 NOM X=[pasta salads], Y=[indian])
Mid2: (Top4 NOM X=[salads], Y=[indian]) -> ∅
```

Nomination continues.

```
Top4: (Top3 NOM X=[salads], Y=[indian]) -> (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0)
```

Top4,
who already “accepted” the Indian food nomination,
now sees that Top3 also accepts it.
Top4 earlier saw Top2 accept it as well.
Top2-Top3-Top4 is one of Top4’s quorums,
and when a quorum accepts something, that’s called _confirmation_.

Top4 confirms the nomination of Indian food and so is ready to begin _balloting_.
A ballot is a <counter,value> pair,
and balloting is the process of finding a ballot that all nodes can commit to.
This happens through multiple rounds of voting on statements about ballots,
beginning with ruling out ballots all nodes can agree _not_ to commit to — so-called “aborted” ballots.

Top4 votes to _prepare_ the ballot <1,indian>.
This means that all lesser ballots are aborted and Top4 promises never to commit to them.
(There are no lesser ballots at this stage,
but anyway that’s the meaning of a “prepare” vote.)

```
Low2: (Top1 NOM X=[burritos indian], Y=[]) -> ∅
Low1: (Top2 NOM X=[], Y=[indian]) -> ∅
```

For Low2 and Low1, nomination continues.

```
Top1: (Top2 NOM X=[], Y=[indian]) -> (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0)
```

Top1 joines Top4 in balloting, also voting to prepare <1,indian>.

```
Top3: (Top4 NOM X=[salads], Y=[indian]) -> ∅
Top2: (Top3 NOM X=[salads], Y=[indian]) -> (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0)
Low1: (Top3 NOM X=[salads], Y=[indian]) -> ∅
Mid4: (Top2 NOM X=[], Y=[indian]) -> ∅
Low1: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Top4: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Mid4: (Top3 NOM X=[salads], Y=[indian]) -> ∅
Mid3: (Top4 NOM X=[salads], Y=[indian]) -> ∅
Low2: (Top4 NOM X=[indian salads], Y=[]) -> ∅
Mid1: (Top3 NOM X=[salads], Y=[indian]) -> ∅
Mid3: (Top2 NOM X=[], Y=[indian]) -> ∅
Mid2: (Top2 NOM X=[], Y=[indian]) -> ∅
Low1: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Top4: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Top4: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Top1: (Top3 NOM X=[salads], Y=[indian]) -> ∅
Mid3: (Top3 NOM X=[salads], Y=[indian]) -> (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0)
Top2: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Top3: (Top2 NOM X=[], Y=[indian]) -> (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0)
Top3: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Low1: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Mid4: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Mid3: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Mid2: (Top3 NOM X=[salads], Y=[indian]) -> (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0)
Mid3: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Low2: (Top4 NOM X=[salads], Y=[indian]) -> ∅
Mid1: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Low1: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top1: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Top3: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Top3: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Mid2: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Top3: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Low1: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top4: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
```

Top4 sees that Top1 is voting to prepare <1,indian>.
This will be important in a moment.

```
Mid2: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Top1: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Top2: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Mid2: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top1: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid4: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Mid3: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid1: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Low2: (Top2 NOM X=[], Y=[indian]) -> ∅
Mid4: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Mid1: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Low1: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
```

Nomination and balloting both continue.
As it happens in this example,
all prepare votes are for the same ballot, <1,indian>.
But it’s possible to have competing prepare votes on differing ballots.

```
Top4: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0)
```

Top4 sees that Top2 is voting to prepare <1,indian>.
He too is voting to prepare <1,indian>,
and he previously saw Top1 vote the same way.
Top1-Top2-Top4 is one of Top4’s quorums,
and a quorum all voting the same way means Top4 can now _accept_ that <1,indian> is prepared.
Top4 sets P to the value of the highest accepted-prepared ballot.

```
Top3: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0)
```

Top3 follows suit upon seeing Top1’s prepare vote.

```
Mid1: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid2: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top2: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Top1: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0)
Low2: (Top3 NOM X=[salads], Y=[indian]) -> ∅
Mid3: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top4: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid1: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Low1: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top3: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Low1: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top2: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid2: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0)
Low1: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Low2: (Top1 NOM X=[burritos], Y=[indian]) -> ∅
Top1: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top3: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top2: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0)
Low2: (Mid2 NOM X=[salads], Y=[indian]) -> ∅
Mid3: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0)
Mid4: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top1: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Low2: (Mid3 NOM X=[pasta salads], Y=[indian]) -> ∅
Top4: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid1: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid2: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top4: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid3: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top3: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top3: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Low2: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid2: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Low1: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top1: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top2: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid4: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid2: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid2: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid1: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top4: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
```

Nomination and balloting both continue.

```
Top3: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
```

Top3 sees that Top1 accepts ballot <1,indian> as prepared.
Top3 does too,
and earlier it saw Top4 accept the same thing.
This makes a quorum all accepting the same ballot,
which means <1,indian> is now _confirmed_ prepared.

Top3 sets CN and HN to the counters in the lowest and highest confirmed-prepared ballots.
Top3 could also in theory have continued to accept new candidates from the nomination phase before now,
but with a confirmed-prepared ballot it no longer can.

By setting CN and HN,
Top3 not only notifies its peers that it confirms <1,indian> is prepared,
but also votes to _commit_ to <1,indian>.

Once a commit vote can be confirmed,
a node considers consensus to be achieved and the value in the ballot can be _externalized_ (acted upon).

```
Low2: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top1: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top2: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid3: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid1: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top4: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Mid2: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Top2: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid1: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Low1: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top3: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Low2: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top1: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Top4: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid3: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Mid1: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid2: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid1: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid1: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Mid1: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top2: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top4: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top1: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Low1: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top3: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Low2: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Top4: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top3: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid3: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid2: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top2: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Top3: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid4: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top3: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid1: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid3: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Mid3: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top4: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low1: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top1: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Low2: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Low1: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top2: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top4: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid4: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid4: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid4: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid4: (Low1 NOM X=[pizza], Y=[]) -> ∅
Mid4: (Low2 NOM X=[sandwiches], Y=[]) -> ∅
Mid4: (Mid1 NOM X=[burgers], Y=[]) -> ∅
Mid4: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid2: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid3: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Low2: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<> PP=<> CN=0 HN=0) -> ∅
Mid1: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid1: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low1: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid3: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
```

More nomination and balloting, more accepting and confirming of prepare votes.

```
Top3: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
```

Top3 has now seen a quorum all confirming preparation of the same ballot. (Top1-Top3-Top4.)
Top3 now _accepts_ that ballot is committed.

```
Low2: (Top4 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Top3: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top1: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top1: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Top3 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid2: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top2: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top4: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Mid3: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top2: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid1: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low1: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Top1 NOM/PREP X=[burritos], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Low1: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Mid2 NOM/PREP X=[salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Mid4: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid1: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top1: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Top3: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Top2 NOM/PREP X=[], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top4: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top2: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid2: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Mid2 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Mid4: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top2: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Top2 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Low1: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid2: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top1: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low1: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid3: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid1: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Mid1 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Top3: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top4: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid1: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid1: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Mid3 NOM/PREP X=[pasta salads], Y=[indian] B=<1,indian> P=<1,indian> PP=<> CN=0 HN=0) -> ∅
Top1: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid2: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top2: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid4: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top4: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top2: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid1: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low1: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Low1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Mid3: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Mid3 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Low2: (Top3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top3: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid4: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Top2: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top2: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid2: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low1: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Low1 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Top4: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top1: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Top4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top4: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid4: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid3: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid1: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid2: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top3: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Top1: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
```

More nodes confirm that <1,indian> is prepared.
Some nodes accept that <1,indian> is committed.
Other nodes are still catching up.

```
Mid4: (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Mid4 EXT C=<1,indian> HN=1)
```

Mid4 has now seen Top1, Top3, and Top4 all accept that <1,indian> is committed.
He is the first to confirm the ballot is committed and so can externalize the value “indian.”

```
Low1: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top4: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Top2: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Mid1: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid3: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low1: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Mid2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top1: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top1: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Top4: (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Top4 EXT C=<1,indian> HN=1)
Top1: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid2: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Top3: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Top1: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Top1 EXT C=<1,indian> HN=1)
Mid3: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Top1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Top2: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Low1: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid3: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid1: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid3: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Top3: (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Top3 EXT C=<1,indian> HN=1)
Mid3: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Top2: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid2: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid1: (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Mid1 EXT C=<1,indian> HN=1)
Mid2: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Mid3: (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Mid3 EXT C=<1,indian> HN=1)
Top2: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Top2 EXT C=<1,indian> HN=1)
Low1: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Low2: (Mid4 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Low2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1)
Low1: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Low1: (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Low1: (Mid2 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Low1 EXT C=<1,indian> HN=1)
Mid2: (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Mid2 EXT C=<1,indian> HN=1)
Low2: (Mid1 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> (Low2 COMMIT B=<1,indian> PN=1 CN=1 HN=1)
Low2: (Top2 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Mid3 PREP B=<1,indian> P=<1,indian> PP=<> CN=1 HN=1) -> ∅
Low2: (Top3 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Low2: (Mid4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Low2: (Top4 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Low2: (Top1 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> ∅
Low2: (Mid2 COMMIT B=<1,indian> PN=1 CN=1 HN=1) -> (Low2 EXT C=<1,indian> HN=1)
```

The rest of the network catches up. By the end, everyone has externalized “indian.”
