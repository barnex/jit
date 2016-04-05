extern void *func_acos;
extern void *func_asin;
extern void *func_atan;
extern void *func_cos;
extern void *func_cosh;
extern void *func_sin;
extern void *func_sinh;
extern void *func_tan;
extern void *func_tanh;
extern void *func_exp;
extern void *func_log;
extern void *func_log10;
extern void *func_sqrt;
extern void *func_fabs;

double eval(void *code, double x, double y);

void eval_2d(void *code, double *dst, double xmin, double xmax, int nx, double ymin, double ymax, int ny);

double call_func(void* f, double x);

