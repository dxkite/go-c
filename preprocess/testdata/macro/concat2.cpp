#define hash_hash # ## #
#define mkstr(a) # a
#define in_between(a) mkstr(a)
#define join(c, d) in_between(c hash_hash d)
char p[] = join(x, y); // equivalent to
// char p[] = "x ## y";
#define concat(a,b,c) a##b##c
#define b BBBB
#define ab AB
concat(Object, a, ccc)
concat(Object, __LINE__, c)
concat(10086 b Object, a, b b )
