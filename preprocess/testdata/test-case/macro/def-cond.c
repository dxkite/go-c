#define mkstr(a) #  a
int main() {
#if 10 + 20 > 30
print(mkstr(10 + 20 > 30));
#else
print(mkstr(not 10+20));
#if 'a'
print(mkstr('a'));
#endif
#endif
}
