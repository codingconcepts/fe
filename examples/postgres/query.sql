SELECT
  pp.proname AS function_name,
  pl.lanname AS function_language,
  pt.typname AS function_return_type,
  pp.proargnames AS function_argument_names,
  ARRAY (
    SELECT pt.typname
    FROM ROWS FROM (unnest(pp.proargtypes))
    WITH ORDINALITY AS a (arg_id, ord)
    JOIN pg_type AS pt ON pt.oid = a.arg_id
    ORDER BY a.ord
  ) AS function_argument_types,
  pp.prosrc AS function_sql
FROM
  pg_proc AS pp
  INNER JOIN pg_namespace AS pn ON pp.pronamespace = pn.oid
  INNER JOIN pg_language AS pl ON pp.prolang = pl.oid
  INNER JOIN pg_type AS pt ON pp.prorettype = pt.oid
WHERE
  pl.lanname NOT IN ('c', 'internal')
  AND pn.nspname NOT LIKE 'pg_%'
  AND pn.nspname != 'information_schema';